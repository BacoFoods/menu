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

var mapLocalesNombres = map[string]string{
	"ZonaG - PC1":      "bacuzonagc14",
	"ZonaG - PC2":      "bacuzonag",
	"Flormorado - PC1": "bacuflormoradopc2",
	"Flormorado - PC2": "bacuflormorado",
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

type Order struct {
	Total struct {
		TotalItems float64 `firestore:"totalItems"`
		Total      float64 `firestore:"total"`
		Propina    float64 `firestore:"propina"`
	} `firestore:"total"`
	Origins  string `firestore:"origen"`
	Payments []struct {
		Total     float64 `firestore:"total"`
		FormaPago string  `firestore:"formaPago"`
	} `firestore:"payments"`
	Pagado        bool      `firestore:"pagado"`
	FechaCreacion time.Time `firestore:"fechaCreacion"`
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
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /temporal/arqueo [get]
func (h *Handler) GetArqueo(c *gin.Context) {
	ctx := context.Background()
	keyLocal := c.Query("local")
	date := c.Query("date") // format 2023-12-30

	minDate := fmt.Sprintf("%s 00:00:00", date)
	maxDate := fmt.Sprintf("%s 23:59:59", date)

	firestore, err := internal.NewFirestore(internal.Config.FirestoreConfig)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	db := make(map[string]map[string]any)

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
	data["orders"] = len(orders)

	var neto1, propinas, bruto float64

	for _, order := range orders {
		neto1 += order.Total.Total
		propinas += order.Total.Propina
		bruto += order.Total.TotalItems

		porPlataforma := make(map[string]float64)
		porMetodo := make(map[string]float64)
		propinaPorMetodo := make(map[string]float64)

		if v, ok := porPlataforma[order.Origins]; ok {
			porPlataforma[order.Origins] = v + order.Total.Total
		} else {
			porPlataforma[order.Origins] = order.Total.Total
		}

		for _, payment := range order.Payments {
			if v, ok := porMetodo[payment.FormaPago]; ok {
				porMetodo[payment.FormaPago] = v + payment.Total
			} else {
				porMetodo[payment.FormaPago] = payment.Total
			}

			propina := 0.0
			if order.Total.Propina > 0 {
				propina = math.Round((0.1 / 1.08) * payment.Total)
			}

			if v, ok := propinaPorMetodo[payment.FormaPago]; ok {
				propinaPorMetodo[payment.FormaPago] = v + propina
			} else {
				propinaPorMetodo[payment.FormaPago] = propina
			}
		}

		data["neto_1"] = neto1
		data["propinas"] = propinas
		data["bruto"] = bruto
		data["descuentos"] = bruto - neto1
		data["ingresos"] = neto1 + propinas
		data["porPlataforma"] = porPlataforma
		data["porMetodoPago"] = porMetodo
		data["propinaPorMetodoPago"] = propinaPorMetodo
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(db))
}
