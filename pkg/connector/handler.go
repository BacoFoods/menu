package connector

import (
	"fmt"
	"net/http"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const LogHandler string = "pkg/connector/handler"

type Handler struct {
	service Service
}

type RequestExcelCreate struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	StoreID   uint   `json:"store_id"`
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
		shared.LogError("error creating equivalence", LogHandler, "Create", err, equivalence)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEquivalenceCreating))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(equivalence))
}

// Find to handle a request to find all equivalence
// @Tags Connector
// @Summary To find equivalence
// @Description To find equivalence
// @Param channel_id query string false "channel_id"
// @Param product_id query string false "product_id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Equivalence}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /equivalence [get]
func (h Handler) Find(c *gin.Context) {
	filters := make(map[string]string)
	if channelID := c.Query("channel_id"); channelID != "" {
		filters["channel_id"] = channelID
	}
	if productID := c.Query("product_id"); productID != "" {
		filters["product_id"] = productID
	}

	equivalences, err := h.service.Find(filters)
	if err != nil {
		shared.LogError("error finding equivalence", LogHandler, "Find", err, equivalences)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEquivalenceFinding))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(equivalences))
}

// Delete to handle a request to delete an equivalence
// @Tags Connector
// @Summary To delete an equivalence
// @Description To delete an equivalence
// @Param id path string true "equivalenceID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Equivalence}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /equivalence/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	equivalenceID := ctx.Param("id")
	equivalence, err := h.service.Delete(equivalenceID)
	if err != nil {
		shared.LogError("error deleting equivalence", LogHandler, "Delete", err, equivalenceID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEquivalenceDeleting))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(equivalence))
}

// Update to handle a request to update an equivalence
// @Tags Connector
// @Summary To update equivalence
// @Description To update equivalence
// @Param equivalence body Equivalence true "equivalence to update"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Equivalence}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /equivalence [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var requestBody Equivalence
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogError("error getting equivalence request body", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorEquivalenceBadRequest))
		return
	}
	equivalence, err := h.service.Update(requestBody)
	if err != nil {
		shared.LogError("error updating equivalence", LogHandler, "Update", err, equivalence)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEquivalenceUpdating))
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
		shared.LogWarn("warning binding request fail", LogHandler, "CreateFile", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	logrus.Infof("creating connector file for store %d [%s, %s]", requestBody.StoreID, requestBody.StartDate, requestBody.EndDate)

	invoices, err := h.service.GetInvoices(requestBody.StartDate, requestBody.EndDate, fmt.Sprint(requestBody.StoreID))
	if err != nil {
		shared.LogError("error getting invoices", LogHandler, "CreateFile", err, invoices)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	if len(invoices) == 0 {
		shared.LogWarn("no invoices found", LogHandler, "CreateFile", nil, invoices)
		ctx.JSON(http.StatusNotFound, shared.ErrorResponse("no invoices for parameters"))
		return
	}

	logrus.Infof("found %d invoices for connector file", len(invoices))

	// Call the HandleSIESAIntegration function to get the Excel file as a byte slice
	excelFile, err := h.service.CreateFile(requestBody.StoreID, invoices)
	if err != nil {
		shared.LogError("error handling SIESA integration", LogHandler, "CreateFile", err, invoices)
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
