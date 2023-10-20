package temporal

import (
	"context"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"time"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

const minDate = "2023-10-07 00:00:00"
const maxDate = "2023-10-07 23:59:59"

var mapLocales = map[string]string{
	"bacuzonagc14":        "ZonaG",
	"bacuzonag":           "ZonaG",
	"bacuflormoradopc2":   "Flormorado",
	"bacuflormorado":      "Flormorado",
	"bacucalle109":        "CL109",
	"bacu109":             "CL109",
	"feriadelmillon2":     "CL90",
	"bacucalle90delivery": "CL90",
	"connectasalon110665": "Connecta",
	"bacuconnecta":        "Connecta",
	"connectasalon210666": "Connecta",
	"cityusalon1":         "CityU",
	"cityusalon2":         "CityU",
	"bacucityu":           "CityU",
	"bacucolinapc110881":  "Colina",
	"bacutitansalon10883": "Titan",
	"bacunogalespc110884": "Nogal",
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

	firestore, err := internal.NewFirestore(internal.Config.FirestoreConfig)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	db := make(map[string]map[string]any)

	for local, nombre := range mapLocales {
		iter := firestore.Client.Collection("confLocalPO/"+local+"/pedidos").
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
			doc.DataTo(&order)
			orders = append(orders, order)
		}

		db[nombre] = make(map[string]any)
		db[nombre]["orders"] = float64(len(orders))

		var neto1, propinas, bruto float64
		for _, order := range orders {
			neto1 += order.Total.Total
			propinas += order.Total.Propina
			bruto += order.Total.TotalItems

			porPlataforma := make(map[string]float64)
			porMetodo := make(map[string]float64)
			propinaPorMetodo := make(map[string]float64)

			for _, payment := range order.Payments {
				porPlataforma[order.Origins] += order.Total.Total
				porMetodo[payment.FormaPago] += payment.Total

				propina := payment.Total
				if order.Total.Propina > 0 {
					propina *= (0.1 / 1.08)
				}
				propinaPorMetodo[payment.FormaPago] += propina
			}

			db[nombre]["neto_1"] = neto1
			db[nombre]["propinas"] = propinas
			db[nombre]["bruto"] = bruto
			db[nombre]["descuentos"] = bruto - neto1
			db[nombre]["ingresos"] = neto1 + propinas
			db[nombre]["porplataforma"] = porPlataforma
			db[nombre]["pormetodo"] = porMetodo
			db[nombre]["propinapormetodo"] = propinaPorMetodo
		}
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(db))
}
