package temporal

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

const (
	LogHandler                     = "pkg/temporal/handler"
	ErrorTemporalArqueoInvalidDate = "invalid date format"
)

var mapLocalesNombres = map[string]string{
	"ZonaG - PC1":      "bacuzonagc14",
	"ZonaG - PC2":      "bacuzonag",
	"Flormorado - PC1": "bacuflormoradopc2",
	"Flormorado - PC2": "bacuflormorado",
	"Flormorado - PC3": "flormorado10885",
	"CL109 - PC1":      "bacucalle109",
	"CL109 - PC2":      "bacu109",
	"CL90 - PC1":       "feriadelmillon2",
	"CL90 - PC2":       "bacucalle90delivery",
	"Connecta - PC1":   "connectasalon110665",
	"Connecta - PC2":   "bacuconnecta",
	"Connecta - PC3":   "connectasalon210666",
	"CityU - PC1":      "cityusalon1",
	"CityU - PC2":      "cityusalon2",
	"CityU - PC3":      "bacucityu",
	"Colina":           "bacucolinapc110881",
	"Titan":            "bacutitansalon10883",
	"Nogal":            "bacunogalespc110884",
}

type OrderPayment struct {
	Total     float64 `firestore:"total"`
	FormaPago string  `firestore:"formaPago"`
}

type Order struct {
	Total struct {
		TotalItems float64 `firestore:"totalItems"`
		Total      float64 `firestore:"total"`
		Propina    float64 `firestore:"propina"`
	} `firestore:"total"`
	Origins       string         `firestore:"origen"`
	Plataforma    string         `firestore:"plataforma"`
	Payments      []OrderPayment `firestore:"payments"`
	Pagado        bool           `firestore:"pagado"`
	FechaCreacion string         `firestore:"fechaCreacion"`
}

// GetLocales to handle the request to get the locales for arqueo
// @Summary Get the locales for arqueo
// @Description Get the locales for arqueo
// @Tags Temporal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /temporal/locales [get]
func (h *Handler) GetLocales(c *gin.Context) {
	c.JSON(http.StatusOK, shared.SuccessResponse(mapLocalesNombres))
}

// GetArqueo to handle the request to get the arqueo
// @Summary Get the arqueo
// @Description Get the arqueo
// @Tags Temporal
// @Accept json
// @Produce json
// @Param local query string true "Local"
// @Param date query string true "Date"
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /temporal/arqueo [get]
func (h *Handler) GetArqueo(c *gin.Context) {
	ctx := context.Background()
	keyLocal := c.Query("local")
	date := c.Query("date") // format 2023-12-30

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		shared.LogWarn("Invalid date format", LogHandler, "GetArqueo", fmt.Errorf(ErrorTemporalArqueoInvalidDate), date)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorTemporalArqueoInvalidDate))
		return
	}

	minDate := fmt.Sprintf("%s 00:00:00", date)
	maxDate := fmt.Sprintf("%s 23:59:59", date)

	firestore, err := internal.NewFirestore(internal.FirestoreConfig(internal.Config.PopappConfig))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	path := fmt.Sprintf("confLocalPO/%s/pedidos", keyLocal)
	iter := firestore.Client.Collection(path).
		Where("pagado", "==", true).
		Where("fechaCreacion", ">=", minDate).
		Where("fechaCreacion", "<=", maxDate).
		Documents(ctx)

	var orders []Order
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error al iterar documentos: %v", err)
		}
		var order Order
		err = doc.DataTo(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		orders = append(orders, order)
	}

	data := make(map[string]any)
	data["Ordenes"] = len(orders)

	var neto1, propinas, bruto float64
	porPlataforma := make(map[string]float64)
	porMetodo := make(map[string]float64)
	propinaPorMetodo := make(map[string]float64)

	for _, order := range orders {
		neto1 += order.Total.Total
		propinas += order.Total.Propina
		bruto += order.Total.TotalItems

		porPlataforma[order.Origins] = porPlataforma[order.Origins] + order.Total.Total

		orderPayments := order.Payments
		if len(orderPayments) == 0 {
			orderPayments = []OrderPayment{
				{
					Total:     order.Total.Total,
					FormaPago: order.Plataforma,
				},
			}
		}

		for _, payment := range orderPayments {
			porMetodo[payment.FormaPago] = porMetodo[payment.FormaPago] + payment.Total

			propinaMP := 0.0
			if order.Total.Propina > 0 {
				propinaMP = math.Round((0.1 / 1.08) * payment.Total)
			}

			propinaPorMetodo[payment.FormaPago] = propinaPorMetodo[payment.FormaPago] + propinaMP
		}

		data["Neto 1"] = neto1
		data["Total propinas"] = propinas
		data["Bruto"] = bruto
		data["Descuentos"] = bruto - neto1
		data["Neto + Propinas"] = neto1 + propinas
		data["Por Plataforma"] = porPlataforma
		data["Por Metodo de Pago"] = porMetodo
		data["Propinas Por Metodo de Pago"] = propinaPorMetodo
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(data))
}
