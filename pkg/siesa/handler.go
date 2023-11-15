package siesa

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler = "pkg/siesa/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

type RequestExcelCreate struct {
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	LocationIDs []string `json:"location_ids"`
}

// Create handles a request to generate an Excel file with orders.
// @Summary Generate an Excel file with orders
// @Description Handles a request to create an Excel file containing orders based on the specified parameters.
// @Tags SIESA
// @Accept json
// @Produce application/vnd.ms-excel
// @Param siesa body RequestExcelCreate true "Request body for creating Excel file"
// @Security ApiKeyAuth
// @Success 200 {file} application/vnd.openxmlformats-officedocument.spreadsheetml.sheet "Successful response with Excel file"
// @Failure 400 {object} shared.Response "Bad Request response"
// @Failure 422 {object} shared.Response "Unprocessable Entity response"
// @Failure 500 {object} shared.Response "Internal Server Error response"
// @Router /siesa [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody RequestExcelCreate
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	response, err := GetOrders(requestBody.StartDate, requestBody.StartDate, requestBody.LocationIDs)
	if err != nil {
		fmt.Println("Error al obtener las ordenes:", err)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	var orders []PopappOrder
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		fmt.Println("Error al decodificar la respuesta JSON:", err)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	// Call the HandleSIESAIntegration function to get the Excel file as a byte slice
	excelFile, err := h.service.HandleSIESAIntegration(orders)
	if err != nil {
		fmt.Println("Error en la integración con SIESA:", err)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=output.xlsx")
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")

	// Send the Excel file as the response
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile)
}

// country, err := h.service.Create(&requestBody)
// if err != nil {
// 	shared.LogError("error creating country", LogHandler, "Create", err, country)
// 	ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingCountry))
// 	return
// }
// ctx.JSON(http.StatusOK, shared.SuccessResponse(country))

// startDate := "2023-10-12"
// endDate := "2023-10-13"
// locationIDs := []string{"bacuzonag"}

// // Obtener las ordenes
// response, err := GetOrders(startDate, endDate, locationIDs)
// if err != nil {
// 	fmt.Println("Error al obtener las ordenes:", err)
// 	return
// }
//fmt.Println("Respuesta del servidor:", response)

// Llamar a la función para manejar la integración con SIESA
