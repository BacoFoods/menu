package siesa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/xuri/excelize/v2"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository}
}

func (s Service) HandleSIESAIntegration(orders []PopappOrder) ([]byte, error) {
	// Construir el cuerpo del request para enviar la información de la orden
	s.insertData()
	doc := s.buildDocument(orders)

	// Generate the Excel file as a byte slice
	file, err := GenerateExcelFile(doc)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	// Save the Excel file to a buffer
	buffer := &bytes.Buffer{}
	if err := file.Write(buffer); err != nil {
		fmt.Println("Error writing Excel file to buffer:", err)
		return nil, err
	}

	// Return the Excel file as a byte slice
	return buffer.Bytes(), nil
}

// getF350IDCO obtiene el IDCO correspondiente a la tienda de la orden. Es decir, el centro de operaciones en SIESA
func getF350IDCO(idStore string) string {
	switch idStore {
	case "bacucityu", "cityusalon1", "cityusalon2":
		return "402"
	case "bacuconnecta", "connectasalon110665", "connectasalon210666":
		return "401"
	case "bacuzonag", "bacuzonagc14":
		return "300"
	case "bacuflormorado", "bacuflormoradopc2":
		return "400"
	case "feriadelmillon2", "bacucalle90delivery":
		return "301"
	case "bacucalle109", "bacu109":
		return "302"
	case "bacudk140":
		return "202"
	case "bacuferia":
		return "301"
	case "bacucolinapc110881":
		return "405"
	case "bacutitansalon10883":
		return "403"
	case "bacunogalespc110884":
		return "303"
	default:
		return "" // Valor predeterminado o manejo de error si es necesario
	}
}

// getF461IDBodegaComponProceso obtiene el ID de la bodega de la tienda de la orden. Es decir, la bodega en SIESA
func getF461IDBodegaComponProceso(idStore string) string {
	switch idStore {
	case "bacucityu", "cityusalon1", "cityusalon2":
		return "402"
	case "bacuconnecta", "connectasalon110665", "connectasalon210666":
		return "401"
	case "bacuzonag", "bacuzonagc14":
		return "300"
	case "bacuflormorado", "bacuflormoradopc2":
		return "400"
	case "feriadelmillon2", "bacucalle90delivery":
		return "301"
	case "bacucalle109", "bacu109":
		return "32A"
	case "bacudk140":
		return "202"
	case "bacuferia":
		return "31E"
	default:
		return "" // Valor predeterminado o manejo de error si es necesario
	}
}

// formatFecha formatea la fecha de la orden al formato requerido por SIESA.
func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		// Manejar el error si ocurre
		return "error en fecha" // Valor por defecto en caso de error
	}
	return t.Format("20060102")
}

// Función que calcula el valor bruto unidades por precio unitario / 1.08
func calculateGrossValue(cantidad int, precioUnitario int) string {
	if precioUnitario == 0 {
		precioUnitario = 1
	}
	precio := precioUnitario
	return strconv.Itoa(precio * cantidad)
}

// buildDocument construye el documento que se enviará al endpoint de Siesa.
func (s Service) buildDocument(orders []PopappOrder) map[string]interface{} {
	doc := make(map[string]interface{})

	doctoVentasComercial := []map[string]string{
		{
			"F350_ID_CO":                    getF350IDCO(orders[0].KeyLocal),                                                                                         // Asigna el valor correspondiente al centro de operación de la primera orden
			"F350_ID_TIPO_DOCTO":            "FVR",                                                                                                                   // Siempre es FVR
			"F350_CONSEC_DOCTO":             "1",                                                                                                                     // TODO: Falta por validar el consecutivo del documento
			"F350_FECHA":                    formatDate(orders[0].FechaCreacion),                                                                                     // Asigna el valor correspondiente a la fecha de la primera orden
			"f461_id_co_fact":               getF350IDCO(orders[0].KeyLocal),                                                                                         // Asigna el valor correspondiente al centro de operación de la primera orden
			"f461_notas":                    "Orden " + orders[0].DisplayID + " - el " + formatDate(orders[0].FechaCreacion) + " - del pdv " + orders[0].NombreStore, // Nota correspondiente a la primera orden procesada
			"F461_ID_BODEGA_COMPON_PROCESO": getF350IDCO(orders[0].KeyLocal),                                                                                         // Asigna el valor correspondiente a la bodega de la primera orden
		},
	}
	doc["Docto. ventas comercial"] = doctoVentasComercial

	// Construir la sección "Descuentos" del documento
	// Puedes iterar sobre los descuentos de la orden y agregarlos a la estructura del documento
	descuentos := []map[string]string{}
	var descuento float64
	registroDescuento := 1 // Número de registro para el descuento
	for _, order := range orders {
		// Obtener el total de items en la orden
		totalItems := float64(order.Total.TotalItems)

		costoEnvio := 0
		if order.Total.CostoEnvio != nil {
			costoEnvio = *order.Total.CostoEnvio
		}

		if costoEnvio != 0 {
			descuento = (totalItems + float64(costoEnvio)) - float64(order.Total.Total)
		} else {
			descuento = totalItems - float64(order.Total.Total)
		}

		if descuento < 0 {
			descuento = 0
		}

		if descuento > float64(order.Total.Total) {
			descuento = float64(order.Total.Total)
		}

		shareDescuento := descuento / totalItems

		for _, item := range order.Items {
			// Ignora si el precio unitario es 0
			// if item.Producto.PrecioUnitario == 0 {
			// 	registroDescuento++ // Incrementar el número de registro sin agregar el descuento
			// 	continue
			// }

			descuentoRegistro := shareDescuento * float64(item.Producto.PrecioUnitario)
			descuentoTotalRegistro := shareDescuento * float64(item.Cantidad) * float64(item.Producto.PrecioUnitario)
			if descuentoRegistro != 0 || descuentoTotalRegistro != 0 {
				descuentoMap := map[string]string{
					"f471_id_co":         getF350IDCO(order.KeyLocal),
					"f471_id_tipo_docto": "FVR",
					"f471_consec_docto":  "1",
					"f471_nro_registro":  strconv.Itoa(registroDescuento),
					"f471_vlr_uni":       strconv.FormatFloat(descuentoRegistro, 'f', 0, 64),      // Formato sin decimales
					"f471_vlr_tot":       strconv.FormatFloat(descuentoTotalRegistro, 'f', 0, 64), // Formato sin decimales
				}
				descuentos = append(descuentos, descuentoMap)
			}
			registroDescuento++

			// Recorrer los itemGroups del item actual
			for _, itemGroup := range item.ItemGroups {
				// Recorrer los modifiers del itemGroup actual
				for _, modifier := range itemGroup.Modifiers {
					// Ignore if unit price is 0
					// if modifier.Producto.PrecioUnitario == 0 {
					// 	registroDescuento++ // Incrementar el número de registro sin agregar el descuento
					// 	continue
					// }

					descuentoRegistro := shareDescuento * (float64(modifier.Producto.PrecioUnitario) / 1.08)
					descuentoTotalRegistro := shareDescuento * (float64(item.Cantidad) * (float64(modifier.Producto.PrecioUnitario) / 1.08))
					if descuentoRegistro != 0 || descuentoTotalRegistro != 0 {
						descuentoMap := map[string]string{
							"f471_id_co":         getF350IDCO(order.KeyLocal),
							"f471_id_tipo_docto": "FVR",
							"f471_consec_docto":  "1",
							"f471_nro_registro":  strconv.Itoa(registroDescuento),
							"f471_vlr_uni":       strconv.FormatFloat(descuentoRegistro, 'f', 0, 64),      // Formato sin decimales
							"f471_vlr_tot":       strconv.FormatFloat(descuentoTotalRegistro, 'f', 0, 64), // Formato sin decimales
						}
						descuentos = append(descuentos, descuentoMap)
					}
					registroDescuento++
				}
			}
		}
	}

	doc["Descuentos"] = descuentos

	// Construir la sección "Movimientos" del documento
	movimientos := []map[string]string{}
	registro := 1 // Variable para el número de registro

	for _, order := range orders {

		for _, item := range order.Items {
			if isValidProduct(item.Producto.Nombre) {
				fmt.Println(order.Tipo, order.Plataforma, item.Producto.Nombre)
				itemMovimiento := map[string]string{
					"f470_id_co":           getF350IDCO(order.KeyLocal),                                         // Asigna el valor correspondiente al centro de operación
					"f470_consec_docto":    "1",                                                                 // Consecutivo del documento auto-incremental
					"f470_nro_registro":    strconv.Itoa(registro),                                              // Asigna el valor correspondiente número de registro cada línea es un producto de la orden
					"f470_id_bodega":       getF461IDBodegaComponProceso(order.KeyLocal),                        // Asigna el valor correspondiente de la bodega
					"f470_id_co_movto":     getF350IDCO(order.KeyLocal),                                         // Asigna el valor correspondiente al centro de operación
					"f470_cant_base":       strconv.Itoa(item.Cantidad),                                         // Asigna la cantidad del item
					"f470_vlr_bruto":       calculateGrossValue(item.Cantidad, item.Producto.PrecioUnitario),    // Valor bruto del item
					"f470_referencia_item": s.GetReferences(order.Tipo, order.Plataforma, item.Producto.Nombre), // TODO: Falta por validar como se hará el cruce de referencias
				}
				movimientos = append(movimientos, itemMovimiento)
				registro++ // Incrementar el número de registro
			}
			// Movement for the item

			// Movements for modifiers
			for _, itemGroup := range item.ItemGroups {
				for _, modifier := range itemGroup.Modifiers {
					if isValidProduct(modifier.Producto.Nombre) {
						fmt.Println(order.Tipo, order.Plataforma, modifier.Producto.Nombre)
						modifierMovimiento := map[string]string{
							"f470_id_co":           getF350IDCO(order.KeyLocal),                                              // Asigna el valor correspondiente al centro de operación
							"f470_consec_docto":    "1",                                                                      // Consecutivo del documento auto-incremental
							"f470_nro_registro":    strconv.Itoa(registro),                                                   // Asigna el valor correspondiente número de registro cada línea es un producto de la orden
							"f470_id_bodega":       getF461IDBodegaComponProceso(order.KeyLocal),                             // Asigna el valor correspondiente de la bodega
							"f470_id_co_movto":     getF350IDCO(order.KeyLocal),                                              // Asigna el valor correspondiente al centro de operación
							"f470_cant_base":       strconv.Itoa(modifier.Cantidad),                                          // Asigna la cantidad del modifier
							"f470_vlr_bruto":       calculateGrossValue(modifier.Cantidad, modifier.Producto.PrecioUnitario), // Asigna el valor correspondiente del modifier
							"f470_referencia_item": s.GetReferences(order.Tipo, order.Plataforma, modifier.Producto.Nombre),  // TODO: Falta por validar como se hará el cruce de referencias
						}
						movimientos = append(movimientos, modifierMovimiento)
						registro++ // Incrementar el número de registro
					}
				}
			}
		}
	}

	doc["Movimientos"] = movimientos

	// Construir la sección "Cuotas CxC" del documento
	cuotasCxC := []map[string]string{}

	itemcuotasCxC := map[string]string{
		"F350_ID_CO":        getF350IDCO(orders[0].KeyLocal),     // Asigna el valor correspondiente al centro de operación
		"F350_CONSEC_DOCTO": "1",                                 //Consecutivo del documento auto-incremental
		"F353_FECHA_VCTO":   formatDate(orders[0].FechaCreacion), // Asigna el valor correspondiente de fecha
	}
	cuotasCxC = append(cuotasCxC, itemcuotasCxC)
	doc["Cuotas CxC"] = cuotasCxC

	//fmt.Println(doc)

	//jsonDoc, err := json.Marshal(doc)
	// if err != nil {
	// 	// Manejar el error si ocurre
	// 	fmt.Printf("Error al convertir a JSON: %v", err)
	// 	return nil
	// }
	// jsonString := string(jsonDoc)
	// fmt.Println(jsonString)
	fmt.Println("--------------------------------------------------")
	return doc
}

// isValidProduct checks if the given product name is valid.
func isValidProduct(productoNombre string) bool {
	invalidProducts := []string{
		"TIENE ENTRADA",
		"SALE FUERTE",
		"Bono 10.000",
		"Bono 20.000",
		"Bono 30.000",
		"Bono 50.000",
		"Bono 75.000",
		"Bono 100.000",
		"En Agua",
		"Sin Azucar",
		"Sin Azúcar",
		"American burguer",
		"Pan negro",
		"Agua",
		"COMBO CHEDDAR",
		"PONLINE",
		"Sin Cuchara, Sin Servilleta",
		"Sin Cuchara",
		"Sin Servilleta",
		"Waffle Pan de Yuca + Jugo Tropical",
		"Tostada Revuelta + Fruta",
		"Pollo Masala + Limonada Natural",
		"Crema Felicidad + Puré de Papa",
		"Huevos Habibi + Limonada Natural",
		"Sin Cuchara",
		"Sin cubiertos",
		"Sin cuchara",
		"En Agua",
		"Sin Azucar",
		"Sin Azúcar",
		"Donación 10k",
		"Donación 20k",
		"Donación 30k",
		"Donación 50k",
		"Donación 75k",
		"Donación 100k",
		"Sin azúcar",
		"Pan negro",
		"American burguer",
		"Agua",
		"COMBO CHEDDAR",
		"PONLINE",
		"Sin Cuchara, Sin Servilleta",
		"Waffle Pan de Yuca + Jugo Tropical",
		"Tostada Revuelta + Fruta",
		"Pollo Masala + Limonada Natural",
		"Crema Felicidad + Puré de Papa",
		"Huevos Habibi + Limonada Natural",
		"Sin Cuchara",
	}

	for _, invalidProduct := range invalidProducts {
		if productoNombre == invalidProduct {
			return false
		}
	}

	return true
}

type OrderResponse struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Orders []PopappOrder `json:"orders"`
}

func GetOrders(startDate string, endDate string, locationIDs []string) (string, error) {
	// Construir la URL
	var allOrders []PopappOrder
	for _, locationID := range locationIDs {
		url := fmt.Sprintf("https://api.popapp.io/orders/?start_date=%s&end_date=%s&location_id=%s&order=desc", startDate, endDate, locationID)
		// Crear la solicitud HTTP
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("error al crear la solicitud: %v", err)
		}
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJiU0ZDUHlqb3NZZkdJQXR1aDYyVUdQcGZvbnUyIiwiY291bnRUb2tlbiI6MSwiaWF0IjoxNjY2NjI2MDcyfQ.EXN9ElWrgtXVRbxLrnBGfPoke_JxMVfgjT7sb6ywdxA")

		// Crear el cliente HTTP
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("error al realizar la solicitud: %v", err)
		}
		defer resp.Body.Close()

		// Leer el cuerpo de la respuesta
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
		}
		// Decodificar la respuesta JSON en la estructura definida
		var orderResponse OrderResponse
		err = json.Unmarshal(body, &orderResponse)
		if err != nil {
			return "", fmt.Errorf("error al decodificar la respuesta JSON: %v", err)
		}

		allOrders = append(allOrders, orderResponse.Orders...)
	}

	// Codificar la lista de órdenes en formato JSON con sangría para mayor legibilidad
	ordersJSON, err := json.MarshalIndent(allOrders, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al codificar la lista de órdenes en formato JSON: %v", err)
	}

	return string(ordersJSON), nil
}

// ExcelColumnID genera el identificador de la columna en formato alfabético
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

// ReferenceError is a custom error type for reference-related errors.
type ReferenceError string

func (e ReferenceError) Error() string {
	return string(e)
}

// GetReferences retrieves references from the database based on order type, platform, and product name.
func (s Service) GetReferences(orderType, platform, productName string) string {
	//var reference Reference
	var filter = make(map[string]string)
	filter["popapp"] = productName
	//var query string
	switch platform {
	case "Popapp":
		switch orderType {
		case "PICK_UP":
			filter["popapp"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DINE_IN":
			// query = `
			// 	SELECT referencia_pdv FROM equivalencias WHERE popapp::text = $1::text LIMIT 1
			// `
			filter["popapp"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaPdv
		case "DELIVERY_BY_RESTAURANT":
			// query = `
			// 	SELECT referencia_pdv FROM equivalencias WHERE popapp::text = $1::text LIMIT 1
			// `
			filter["popapp"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaPdv
		default:
			return convertError(ReferenceError(fmt.Sprintf("unsupported order type: %s for platform: %s", orderType, platform)))
		}
	case "RAPPI":
		switch orderType {
		case "PICK_UP":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE rappi_pick_up::text = $1::text LIMIT 1
			// `
			filter["rappi_pick_up"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_PLATAFORMA":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE rappi_bacu::text = $1::text LIMIT 1
			// `
			filter["rappi_bacu"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_RESTAURANT":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE rappi_bacu::text = $1::text LIMIT 1
			// `
			filter["rappi_bacu"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		default:
			return convertError(ReferenceError(fmt.Sprintf("unsupported order type: %s for platform: %s", orderType, platform)))
		}
	case "DiDi":
		switch orderType {
		case "PICK_UP":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE didi_bacu::text = $1::text LIMIT 1
			// `
			filter["didi_bacu"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_PLATAFORMA":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE didi_bacu::text = $1::text LIMIT 1
			// `
			filter["didi_bacu"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_RESTAURANT":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE didi_bacu::text = $1::text LIMIT 1
			// `
			filter["didi_bacu"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		default:
			return convertError(ReferenceError(fmt.Sprintf("unsupported order type: %s for platform: %s", orderType, platform)))
		}
	case "BACOMARKETPLACE":
		switch orderType {
		case "PICK_UP":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE bacu_marketplace::text = $1::text LIMIT 1
			// `
			filter["bacu_marketplace"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_PLATAFORMA":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE bacu_marketplace::text = $1::text LIMIT 1
			// `
			filter["bacu_marketplace"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_RESTAURANT":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE bacu_marketplace::text = $1::text LIMIT 1
			// `
			filter["bacu_marketplace"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		default:
			return convertError(ReferenceError(fmt.Sprintf("unsupported order type: %s for platform: %s", orderType, platform)))
		}
	case "ORDERINTABLE":
		switch orderType {
		case "PICK_UP":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE bacu_marketplace::text = $1::text LIMIT 1
			// `
			filter["bacu_marketplace"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_PLATAFORMA":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE bacu_marketplace::text = $1::text LIMIT 1
			// `
			filter["bacu_marketplace"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		case "DELIVERY_BY_RESTAURANT":
			// query = `
			// 	SELECT referencia_delivery_inline FROM equivalencias WHERE bacu_marketplace::text = $1::text LIMIT 1
			// `
			filter["bacu_marketplace"] = productName
			reference, err := s.repository.Find(filter)
			if err != nil {
				shared.LogError("error getting reference row", LogDBRepository, "Find", err, filter)
				return " "
			}
			return reference.ReferenciaDeliveryInline
		default:
			return convertError(ReferenceError(fmt.Sprintf("unsupported order type: %s for platform: %s", orderType, platform)))
		}
	default:
		return convertError(ReferenceError(fmt.Sprintf("unsupported platform: %s", platform)))
	}
	//fmt.Println(platform,orderType,productName)
	//row := db.QueryRow(query, productName)

	// err = row.Scan(&reference)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return convertError(ReferenceError("no references found for the given conditions"))
	// 	}
	// 	return convertError(ReferenceError(fmt.Sprintf("error scanning row: %v", err)))
	// }
	return ""
}

// convertError converts a ReferenceError to a string.
func convertError(err ReferenceError) string {
	return string(err)
}

func getElement(row []string, index int) string {
	if index >= 0 && index < len(row) {
		return row[index]
	}
	return "NULL"
}

func (s Service) insertData() error {

	// Abrir el archivo Excel
	xlsx, err := excelize.OpenFile("equivalencias.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// Leer las celdas del archivo Excel y realizar inserciones en la base de datos
	rows, err := xlsx.GetRows("Hoja1") // Asegúrate de que el nombre de la hoja sea correcto
	if err != nil {
		log.Fatal(err)
	}

	// Comenzar desde la segunda fila para evitar la fila de encabezados
	for _, row := range rows[1:] {
		reference := Reference{
			Category:                 getElement(row, 0),
			Popapp:                   getElement(row, 1),
			ReferenciaPdv:            getElement(row, 2),
			ReferenciaDeliveryInline: getElement(row, 3),
			RappiPickUp:              getElement(row, 4),
			RappiBacu:                getElement(row, 5),
			DidiStu:                  getElement(row, 6),
			DidiBacu:                 getElement(row, 7),
			BacuMarketplace:          getElement(row, 8),
		}
		err2 := s.repository.TruncateRecords()
		if err2 != nil {
			shared.LogError(err.Error(), LogDBRepository, "TruncateRecords", err, reference)
		}
		err := s.repository.Create(&reference)
		if err != nil {
			shared.LogError(err.Error(), LogDBRepository, "Create", err, reference)
		}
	}
	fmt.Println("Datos insertados correctamente en la tabla equivalencias.")

	return nil
}

func (s Service) Create(reference *Reference) error {
	return s.repository.Create(reference)
}

func (s Service) Find(query map[string]string) (*Reference, error) {
	return s.repository.Find(query)
}

func (s Service) TruncateRecords() error {
	return s.repository.TruncateRecords()
}
