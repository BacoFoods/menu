package connector

import (
	"bytes"
	"errors"
	"fmt"
	invoicePkg "github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/shared"
	storePkg "github.com/BacoFoods/menu/pkg/store"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

const (
	LogService string = "pkg/connector/service"
)

type service struct {
	repository Repository
	invoice    invoicePkg.Repository
	store      storePkg.Repository
}

func NewService(repository *DBRepository, invoice invoicePkg.Repository, store storePkg.Repository) service {
	return service{repository, invoice, store}
}

func (s service) Create(equivalence *Equivalence) (*Equivalence, error) {
	return s.repository.Create(equivalence)
}

func (s service) Find(filter map[string]string) ([]Equivalence, error) {
	return s.repository.Find(filter)
}

func (s service) Update(equivalence Equivalence) (*Equivalence, error) {
	return s.repository.Update(equivalence)
}

func (s service) Delete(equivalenceID string) (*Equivalence, error) {
	return s.repository.Delete(equivalenceID)
}

func (s service) GetInvoices(startDate, endDate, storeID string) ([]invoicePkg.Invoice, error) {
	filter := map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
		"store_id":   storeID,
	}

	invoices, err := s.invoice.FindInvoices(filter)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s service) CreateFile(invoices []invoicePkg.Invoice) ([]byte, error) {
	doc := s.BuildDocument(invoices)
	fmt.Println(doc)

	file, err := GenerateExcelFile(doc)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	buffer := &bytes.Buffer{}
	if err := file.Write(buffer); err != nil {
		fmt.Println("Error writing Excel file to buffer:", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func getF350IDCO(storeID string, storeRepo storePkg.Repository) string {
	store, err := storeRepo.Get(storeID)
	if err != nil {
		shared.LogError(fmt.Sprintf("Error getting store: %v", err), LogService, "getF350IDCO", nil, storeID)
		// You might want to return a default or meaningful value here based on your requirements
		return ""
	}
	return store.OpsCenter
}

func getF461IDCO(storeID string, storeRepo storePkg.Repository) string {
	store, err := storeRepo.Get(storeID)
	if err != nil {
		shared.LogError(fmt.Sprintf("Error getting store: %v", err), LogService, "getF350IDCO", nil, storeID)
		// You might want to return a default or meaningful value here based on your requirements
		return ""
	}
	return store.Wharehouse
}

func formatDate(createdAt *time.Time) string {
	if createdAt == nil {
		return "error en fecha"
	}
	return createdAt.Format("20060102")
}

func calculateGrossValue(cantidad float64, precioUnitario float64) string {
	if precioUnitario == 0 {
		precioUnitario = 1
	}
	precio := precioUnitario
	return strconv.FormatFloat(precio*cantidad, 'f', 0, 64)
}

func (s service) GetReferences(channelID, productID string) string {
	filter := map[string]string{
		"channel_id": channelID,
		"product_id": productID,
	}
	equivalence, err := s.repository.FindReference(filter)

	if err != nil {
		shared.LogError("error getting SiesaID", LogService, "GetReferences", err, channelID, productID)
		return ""
	}

	if equivalence != nil {
		return equivalence.SiesaID
	}

	return ""
}

func (s service) BuildDocument(invoices []invoicePkg.Invoice) map[string]interface{} {
	doc := make(map[string]interface{})

	doctoVentasComercial := []map[string]string{
		{
			"F350_ID_CO":                    getF350IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
			"F350_ID_TIPO_DOCTO":            "FVR",
			"F350_CONSEC_DOCTO":             "1",
			"F350_FECHA":                    formatDate(invoices[0].CreatedAt),
			"f461_id_co_fact":               getF350IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
			"f461_notas":                    formatDate(invoices[0].CreatedAt) + " - del pdv " + strconv.Itoa(int(*invoices[0].StoreID)),
			"F461_ID_BODEGA_COMPON_PROCESO": getF461IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
		},
	}
	doc["Docto. ventas comercial"] = doctoVentasComercial

	// Construir la sección "Descuentos" del documento
	descuentos := []map[string]string{}
	registroDescuento := 1 // Número de registro para el descuento

	for _, invoice := range invoices {
		for _, item := range invoice.Items {

			descuentoRegistro := item.Price - item.DiscountedPrice
			descuentoTotalRegistro := item.Price - item.DiscountedPrice

			if descuentoRegistro != 0 || descuentoTotalRegistro != 0 {
				descuentoMap := map[string]string{
					"f471_id_co":         getF350IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
					"f471_id_tipo_docto": "FVR",
					"f471_consec_docto":  "1",
					"f471_nro_registro":  strconv.Itoa(registroDescuento),
					"f471_vlr_uni":       strconv.FormatFloat(descuentoRegistro, 'f', 0, 64),
					"f471_vlr_tot":       strconv.FormatFloat(descuentoTotalRegistro, 'f', 0, 64),
				}
				descuentos = append(descuentos, descuentoMap)
			}
			registroDescuento++
		}
	}

	doc["Descuentos"] = descuentos

	// Construir la sección "Movimientos" del documento
	movimientos := []map[string]string{}
	registro := 1 // Variable para el número de registro

	for _, invoice := range invoices {
		for _, item := range invoice.Items {

			itemMovimiento := map[string]string{
				"f470_id_co":           getF350IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
				"f470_consec_docto":    "1",
				"f470_nro_registro":    strconv.Itoa(registro),
				"f470_id_bodega":       getF461IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
				"f470_id_co_movto":     getF350IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
				"f470_vlr_bruto":       calculateGrossValue(1.0, item.Price),
				"f470_referencia_item": s.GetReferences(strconv.Itoa(int(*invoice.ChannelID)), strconv.Itoa(int(*item.ProductID))),
			}
			movimientos = append(movimientos, itemMovimiento)
			registro++
		}
	}

	doc["Movimientos"] = movimientos

	// Construir la sección "Cuotas CxC" del documento
	cuotasCxC := []map[string]string{}

	itemcuotasCxC := map[string]string{
		"F350_ID_CO":        getF350IDCO(strconv.Itoa(int(*invoices[0].StoreID)), s.store),
		"F350_CONSEC_DOCTO": "1",
		"F353_FECHA_VCTO":   formatDate(invoices[0].CreatedAt),
	}
	cuotasCxC = append(cuotasCxC, itemcuotasCxC)
	doc["Cuotas CxC"] = cuotasCxC

	return doc
}

func ExcelColumnID(colIdx int) string {
	var result string
	for {
		if colIdx > 0 {
			colIdx--
			result = string('A'+colIdx%26) + result
			colIdx /= 26
		} else {
			break
		}
	}
	return result
}

func GenerateExcelFile(doc map[string]interface{}) (*excelize.File, error) {
	file := excelize.NewFile()

	// Delete the default "Sheet1"
	file.DeleteSheet("Sheet1")

	// Define sheet names and corresponding headers
	sheetNames := []string{"Docto. ventas comercial", "Descuentos", "Cuotas CxC", "Movimientos"}
	headersColumns := [][]string{
		{"F350_ID_CO", "F350_ID_TIPO_DOCTO", "F350_CONSEC_DOCTO", "F350_FECHA", "f461_id_co_fact", "f461_notas", "F461_ID_BODEGA_COMPON_PROCESO"},
		{"f471_id_co", "f471_id_tipo_docto", "f471_consec_docto", "f471_nro_registro", "f471_vlr_uni", "f471_vlr_tot"},
		{"F350_ID_CO", "F350_CONSEC_DOCTO"},
		{"f470_id_co", "f470_consec_docto", "f470_nro_registro", "f470_id_bodega", "f470_id_co_movto", "f470_cant_base", "f470_vlr_bruto", "f470_referencia_item"},
	}

	// Define headers for each sheet
	headers := [][]string{
		{"Centro Operacion", "Tipo Documento", "Numero Docto", "Fecha Docto", "Centro operación factura", "Observaciones", "Bodega componentes Kit"},
		{"Centro Operacion", "Tipo Documento", "Consecutivo Documento", "Numero Registro", "Valor Descuento Unitario", "Valor Descuento Total"},
		{"Centro Operacion", "Número Documento"},
		{"Centro Operacion", "Consecutivo Documento", "Numero Registro", "Bodega", "Centro Operacion Mvmnto", "Cantidad", "Valor Neto", "Referencia"},
	}

	// Create sheets and add headers
	for i, sheetName := range sheetNames {
		index, err := file.NewSheet(sheetName)
		if err != nil {
			return nil, err
		}
		file.SetActiveSheet(index)

		// Add headers to the sheet
		for col, header := range headers[i] {
			cell := ExcelColumnID(col+1) + "1"
			file.SetCellValue(sheetName, cell, header)
		}

		// Extract data from the document and add to the sheet
		if data, ok := doc[sheetName]; ok {
			if records, ok := data.([]map[string]string); ok {
				for rowIdx, record := range records {
					for colIdx, header := range headersColumns[i] {
						value, ok := record[header]
						if !ok {
							return nil, fmt.Errorf("missing value for header %s in sheet %s", headersColumns, sheetName)
						}
						cell := ExcelColumnID(colIdx+1) + fmt.Sprint(rowIdx+2) // Start from row 2 for data
						file.SetCellValue(sheetName, cell, value)
					}
				}
			}
		}
	}
	// Delete the default "Sheet1"
	file.DeleteSheet("Sheet1")

	// Check if the file is empty
	if file.SheetCount == 0 {
		return nil, errors.New("no data found, Excel file is empty")
	}

	return file, nil
}

type Service interface {
	Create(*Equivalence) (*Equivalence, error)
	Find(filter map[string]string) ([]Equivalence, error)
	Update(Equivalence) (*Equivalence, error)
	Delete(string) (*Equivalence, error)
	CreateFile(invoices []invoicePkg.Invoice) ([]byte, error)
	GetInvoices(startDate, endDate, storeID string) ([]invoicePkg.Invoice, error)
}
