package connector

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/connector/handler"

type Handler struct {
	service Service
}

type RequestExcelCreate struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	StoreID   string `json:"store_id"`
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create to handle a request to create an equivalence
// @Tags Connector
// @Summary To create an equivalence
// @Description To create an equivalence
// @Param equivalence body Equivalence true "equivalence request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Equivalence}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /equivalence [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody Equivalence
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorEquivalenceBadRequest))
		return
	}

	equivalence, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating country", LogHandler, "Create", err, equivalence)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEquivalenceCreating))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(equivalence))
}

// Create handles a request to generate an Excel file with orders invoice.
// @Summary Generate an Excel file with orders
// @Description Handles a request to create an Excel file containing orders invoice.
// @Tags Connector
// @Accept json
// @Produce application/vnd.ms-excel
// @Param connector body RequestExcelCreate true "Request body for creating Excel file"
// @Security ApiKeyAuth
// @Success 200 {file} application/vnd.openxmlformats-officedocument.spreadsheetml.sheet "Successful response with Excel file"
// @Failure 400 {object} shared.Response "Bad Request response"
// @Failure 422 {object} shared.Response "Unprocessable Entity response"
// @Failure 500 {object} shared.Response "Internal Server Error response"
// @Router /connector [post]
func (h *Handler) CreateFile(ctx *gin.Context) {
	var requestBody RequestExcelCreate
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	invoices, err := h.service.GetInvoices(requestBody.StartDate, requestBody.EndDate, requestBody.StoreID)
	if err != nil {
		fmt.Println("Error getting invoices:", err)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	// Call the HandleSIESAIntegration function to get the Excel file as a byte slice
	excelFile, err := h.service.CreateFile(invoices)
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
