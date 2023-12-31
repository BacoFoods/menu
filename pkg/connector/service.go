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
	"github.com/sirupsen/logrus"
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

func (s service) GetInvoices(startDate, endDate string, storeID []string) ([]invoicePkg.Invoice, error) {
	filter := map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
		"stores":     storeID,
		"status":     "paid", // TODO: use const
	}

	invoices, err := s.invoice.FindInvoices(filter)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s service) CreateFile(stores []uint, invoices []invoicePkg.Invoice) ([]byte, error) {
	doc, err := s.BuildDocument(stores, invoices)
	if err != nil {
		shared.LogError("error generating building document", LogService, "CreateFile-BuildDocument", err)
		return nil, err
	}

	file, err := GenerateExcelFile(doc)
	if err != nil {
		shared.LogError("error generating Excel file", LogService, "CreateFile-GenerateExcelFile", err)
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

func (s service) GetReferences() ([]Equivalence, error) {
	return s.repository.FindReference()
}

func (s service) BuildDocument(stores []uint, invoices []invoicePkg.Invoice) (map[string]interface{}, error) {
	doc := make(map[string]interface{})
	storesDb := make([]*storePkg.Store, 0)
	for _, storeID := range stores {
		store, err := s.store.Get(fmt.Sprint(storeID))
		if err != nil {
			shared.LogError("error getting store", LogService, "BuildDocument", err, storeID)
			return nil, err
		}

		storesDb = append(storesDb, store)
	}

	if len(storesDb) == 0 {
		return nil, errors.New("no stores found")
	}

	store := storesDb[0]
	// Check all stores share the same OpsCenter and Wharehouse
	for _, s := range storesDb {
		if s.OpsCenter != store.OpsCenter {
			return nil, errors.New("stores don't share the same OpsCenter - " + s.OpsCenter + " - " + store.OpsCenter + " -")
		}
		if s.Wharehouse != store.Wharehouse {
			return nil, errors.New("stores don't share the same Wharehouse - " + s.Wharehouse + " - " + store.Wharehouse + " -")
		}
	}

	f350IDCO := getF350IDCO(store)
	f461IDCO := getF461IDCO(store)
	equivalences, err := s.GetReferences()
	if err != nil {
		fmt.Println("Error al obtener equivalencias:", err)
		return nil, err
	}

	equivalencesMap := make(map[string]string)
	for _, equivalence := range equivalences {
		key := fmt.Sprintf("%s_%s", equivalence.ChannelID, equivalence.ProductID)
		equivalencesMap[key] = equivalence.SiesaID
	}

	now := invoices[0].CreatedAt
	if now == nil {
		now = invoices[0].UpdatedAt
	}

	if now == nil {
		d := time.Now()
		now = &d
	}

	doc["Docto. ventas comercial"] = []map[string]string{
		{
			"F350_ID_CO":         f350IDCO,
			"F350_ID_TIPO_DOCTO": "FVR",
			// TODO: Revisar si se puede cambiar el consecutivo del documento
			"F350_CONSEC_DOCTO":             "1",
			"F350_FECHA":                    formatDate(*now),
			"f461_id_co_fact":               f350IDCO,
			"f461_notas":                    fmt.Sprintf("%s  - del pdv %v", formatDate(*now), stores),
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
			key := fmt.Sprintf("%s_%s", fmt.Sprint(*invoice.ChannelID), fmt.Sprint(*item.ProductID))
			siesaID, exists := equivalencesMap[key]
			if !exists {
				shared.LogError("No se encontró equivalencia", LogService, "BuildDocument", nil, "ChannelID", fmt.Sprint(*invoice.ChannelID), "ProductID", fmt.Sprint(*item.ProductID))
				continue
			}

			itemMovimiento := map[string]string{
				"f470_id_co":           f350IDCO,
				"f470_consec_docto":    "1",
				"f470_nro_registro":    strconv.Itoa(registro),
				"f470_id_bodega":       f461IDCO,
				"f470_id_co_movto":     f350IDCO,
				"f470_cant_base":       "1",
				"f470_vlr_bruto":       calculateGrossValue(1, item.Price),
				"f470_referencia_item": siesaID,
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
		"F353_FECHA_VCTO":   formatDate(*now),
	}
	cuotasCxC = append(cuotasCxC, itemcuotasCxC)
	doc["Cuotas CxC"] = cuotasCxC

	return doc, nil
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

type sheetCol struct {
	Header    string
	FieldName string
}

func GenerateExcelFile(doc map[string]interface{}) (*excelize.File, error) {
	file := excelize.NewFile()

	// Delete the default "Sheet1"
	file.DeleteSheet("Sheet1")

	// Define sheet names and corresponding headers
	sheets := map[string][]sheetCol{
		"Docto. ventas comercial": {
			{
				Header:    "Centro Operacion",
				FieldName: "F350_ID_CO",
			},
			{
				Header:    "Tipo Documento",
				FieldName: "F350_ID_TIPO_DOCTO",
			},
			{
				Header:    "Numero Docto",
				FieldName: "F350_CONSEC_DOCTO",
			},
			{
				Header:    "Fecha Docto",
				FieldName: "F350_FECHA",
			},
			{
				Header:    "Centro operación factura",
				FieldName: "f461_id_co_fact",
			},
			{
				Header:    "Observaciones",
				FieldName: "f461_notas",
			},
			{
				Header:    "Bodega componentes Kit",
				FieldName: "F461_ID_BODEGA_COMPON_PROCESO",
			},
		},
		"Descuentos": {
			{
				Header:    "Centro Operacion",
				FieldName: "f471_id_co",
			},
			{
				Header:    "Tipo Documento",
				FieldName: "f471_id_tipo_docto",
			},
			{
				Header:    "Consecutivo Documento",
				FieldName: "f471_consec_docto",
			},
			{
				Header:    "Numero Registro",
				FieldName: "f471_nro_registro",
			},
			{
				Header:    "Valor Descuento Unitario",
				FieldName: "f471_vlr_uni",
			},
			{
				Header:    "Valor Descuento Total",
				FieldName: "f471_vlr_tot",
			},
		},
		"Cuotas CxC": {
			{
				Header:    "Centro Operacion",
				FieldName: "F350_ID_CO",
			},
			{
				Header:    "Número Documento",
				FieldName: "F350_CONSEC_DOCTO",
			},
		},
		"Movimientos": {
			{
				Header:    "Centro Operacion",
				FieldName: "f470_id_co",
			},
			{
				Header:    "Consecutivo Documento",
				FieldName: "f470_consec_docto",
			},
			{
				Header:    "Numero Registro",
				FieldName: "f470_nro_registro",
			},
			{
				Header:    "Bodega",
				FieldName: "f470_id_bodega",
			},
			{
				Header:    "Centro Operacion Mvmnto",
				FieldName: "f470_id_co_movto",
			},
			{
				Header:    "Cantidad",
				FieldName: "f470_cant_base",
			},
			{
				Header:    "Valor Neto",
				FieldName: "f470_vlr_bruto",
			},
			{
				Header:    "Referencia",
				FieldName: "f470_referencia_item",
			},
		},
	}

	// Create sheets and add headers
	for name, cols := range sheets {
		currentSheedIx, err := file.NewSheet(name)
		if err != nil {
			return nil, err
		}

		file.SetActiveSheet(currentSheedIx)

		// Add headers to the sheet
		for colIdx, col := range cols {
			cell := ExcelColumnID(colIdx+1) + "1"
			file.SetCellValue(name, cell, col.Header)
		}

		// Extract data from the document and add to the sheet
		if data, ok := doc[name]; ok {
			if records, ok := data.([]map[string]string); ok {
				for rowIdx, record := range records {
					for colIdx, header := range cols {
						value, ok := record[header.FieldName]
						if !ok {
							logrus.Warnf("missing value for header %s in sheet %s", header.FieldName, name)
							value = ""
						}

						cell := ExcelColumnID(colIdx+1) + fmt.Sprint(rowIdx+2) // Start from row 2 for data
						_ = file.SetCellValue(name, cell, value)
					}
				}
			}
		}
	}

	// Delete the default "Sheet1"
	_ = file.DeleteSheet("Sheet1")

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
	CreateFile(stores []uint, invoices []invoicePkg.Invoice) ([]byte, error)
	GetInvoices(startDate, endDate string, storeIDs []string) ([]invoicePkg.Invoice, error)
}
