package connector

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	invoicePkg "github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/shared"
	storePkg "github.com/BacoFoods/menu/pkg/store"
	"github.com/xuri/excelize/v2"
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

func (s service) CreateFile(storeID uint, invoices []invoicePkg.Invoice) ([]byte, error) {
	doc, err := s.BuildDocument(storeID, invoices)
	if err != nil {
		return nil, err
	}

	file, err := GenerateExcelFile(doc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := &bytes.Buffer{}
	if err := file.Write(buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func getF350IDCO(s *storePkg.Store) string {
	if s == nil {
		return ""
	}

	return s.OpsCenter
}

func getF461IDCO(s *storePkg.Store) string {
	if s == nil {
		return ""
	}

	return s.Wharehouse
}

func formatDate(t time.Time) string {
	return t.Format("20060102")
}

func calculateGrossValue(cantidad int, precioUnitario float64) string {
	if precioUnitario == 0 {
		precioUnitario = 1
	}

	return strconv.FormatFloat(precioUnitario*float64(cantidad), 'f', 0, 64)
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

func (s service) BuildDocument(storeID uint, invoices []invoicePkg.Invoice) (map[string]interface{}, error) {
	doc := make(map[string]interface{})
	store, err := s.store.Get(fmt.Sprint(storeID))
	if err != nil {
		shared.LogError("error getting store", LogService, "BuildDocument", err, storeID)
		return nil, err
	}

	f350IDCO := getF350IDCO(store)
	f461IDCO := getF461IDCO(store)
	now := time.Now()

	doc["Docto. ventas comercial"] = []map[string]string{
		{
			"F350_ID_CO":         f350IDCO,
			"F350_ID_TIPO_DOCTO": "FVR",
			// TODO: Revisar si se puede cambiar el consecutivo del documento
			"F350_CONSEC_DOCTO":             "1",
			"F350_FECHA":                    formatDate(now),
			"f461_id_co_fact":               f350IDCO,
			"f461_notas":                    fmt.Sprintf("%s  - del pdv %d", formatDate(now), storeID),
			"F461_ID_BODEGA_COMPON_PROCESO": f461IDCO,
		},
	}

	// Construir la sección "Descuentos" del documento
	descuentos := []map[string]string{}
	registroDescuento := 0 // Número de registro para el descuento

	for _, invoice := range invoices {
		for _, item := range invoice.Items {
			descuentoRegistro := item.Price - item.DiscountedPrice

			if descuentoRegistro == 0 {
				continue
			}

			registroDescuento++
			descuentos = append(descuentos, map[string]string{
				"f471_id_co":         f350IDCO,
				"f471_id_tipo_docto": "FVR",
				"f471_consec_docto":  "1",
				"f471_nro_registro":  strconv.Itoa(registroDescuento),
				"f471_vlr_uni":       strconv.FormatFloat(descuentoRegistro, 'f', 0, 64),
				"f471_vlr_tot":       strconv.FormatFloat(descuentoRegistro, 'f', 0, 64),
			})
		}
	}

	doc["Descuentos"] = descuentos

	// Construir la sección "Movimientos" del documento
	movimientos := []map[string]string{}
	registro := 1 // Variable para el número de registro

	for _, invoice := range invoices {
		for _, item := range invoice.Items {
			itemMovimiento := map[string]string{
				"f470_id_co":           f350IDCO,
				"f470_consec_docto":    "1",
				"f470_nro_registro":    strconv.Itoa(registro),
				"f470_id_bodega":       f461IDCO,
				"f470_id_co_movto":     f350IDCO,
				"f470_vlr_bruto":       calculateGrossValue(1, item.Price),
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
		"F350_ID_CO":        f350IDCO,
		"F350_CONSEC_DOCTO": "1",
		"F353_FECHA_VCTO":   formatDate(now),
	}
	cuotasCxC = append(cuotasCxC, itemcuotasCxC)
	doc["Cuotas CxC"] = cuotasCxC

	return doc, err
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
	CreateFile(storeID uint, invoices []invoicePkg.Invoice) ([]byte, error)
	GetInvoices(startDate, endDate, storeID string) ([]invoicePkg.Invoice, error)
}
