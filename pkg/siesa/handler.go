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
	response, err := GetOrders(requestBody.StartDate, requestBody.EndDate, requestBody.LocationIDs)
	//fmt.Println(response)
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

// FindReferences handles a request to find references based on filters
// @Tags SIESA
// @Summary Find references
// @Description Find references based on filters
// @Param referencia_pdv query string false "PDV Reference"
// @Param referencia_delivery_inline query string false "Delivery Inline Reference"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /siesa/reference [get]
func (h *Handler) FindReferences(ctx *gin.Context) {
	filters := make(map[string]string)

	referenciaPDV := ctx.Query("referencia_pdv")
	if referenciaPDV != "" {
		filters["referencia_pdv"] = referenciaPDV
	}

	referenciaDeliveryInline := ctx.Query("referencia_delivery_inline")
	if referenciaDeliveryInline != "" {
		filters["referencia_delivery_inline"] = referenciaDeliveryInline
	}

	// Check if both filters are empty, and if so, don't apply any filters
	if referenciaPDV == "" && referenciaDeliveryInline == "" {
		filters = nil
	}

	references, err := h.service.FindReferences(filters)
	if err != nil {
		shared.LogError("error getting references", LogHandler, "FindReferences", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingReferences))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(references))
}

// CreateReference handles a request to create a reference
// @Tags SIESA
// @Summary Create a reference
// @Description Create a reference
// @Param reference body Reference true "Reference request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /siesa/reference [post]
func (h *Handler) CreateReference(ctx *gin.Context) {
	var requestBody Reference
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request failed", LogHandler, "CreateReference", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	reference, err := h.service.CreateReference(&requestBody)
	if err != nil {
		shared.LogError("error creating reference", LogHandler, "CreateReference", err, reference)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingReference))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(reference))
}

// Delete to handle a request to delete a reference
// @Tags SIESA
// @Summary To delete a reference
// @Description To delete a reference
// @Param id path string true "reference id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Reference}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /siesa/reference/{id} [delete]
func (h *Handler) DeleteReference(c *gin.Context) {
	referenceID := c.Param("id")
	reference, err := h.service.DeleteReference(referenceID)
	if err != nil {
		shared.LogError("error deleting category", LogHandler, "DeleteReference", err, reference)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingReference))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(reference))
}

// Update to handle a request to update a reference
// @Tags SIESA
// @Summary To update a reference
// @Description To update a reference
// @Param id path string true "reference id"
// @Param reference body Reference true "reference request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Reference}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /siesa/reference/{id} [patch]
func (h *Handler) UpdateReference(c *gin.Context) {
	var requestBody Reference
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "UpdateReference", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	category, err := h.service.UpdateReference(&requestBody)
	if err != nil {
		shared.LogError("error updating category", LogHandler, "Update", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingReference))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}
