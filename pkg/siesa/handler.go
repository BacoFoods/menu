package siesa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler = "pkg/siesa/handler"

var mapLocalesNombres = map[string][]string{
	"ZonaG":      {"bacuzonagc14", "bacuzonag"},
	"Flormorado": {"bacuflormoradopc2", "bacuflormorado", "flormorado10885"},
	"CL109":      {"bacucalle109", "bacu109"},
	"CL90":       {"feriadelmillon2", "bacucalle90delivery"},
	"Connecta":   {"connectasalon110665", "bacuconnecta", "connectasalon210666"},
	"CityU":      {"cityusalon1", "cityusalon2", "bacucityu"},
	"Colina":     {"bacucolinapc110881"},
	"Titan":      {"bacutitansalon10883"},
	"Nogal":      {"bacunogalespc110884"},
}

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

// GetLocales to handle the request to get the locales for SIESA
// @Summary Get the locales for SIESA
// @Description Get the locales for SIESA
// @Tags SIESA
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /siesa/locales [get]
func (h *Handler) GetLocales(c *gin.Context) {
	c.JSON(http.StatusOK, shared.SuccessResponse(mapLocalesNombres))
}

// CreateJSON handles a request to generate an Excel file with orders.
// @Summary Generate an JSON with orders
// @Description Handles a request to create an JSON containing orders based on the specified parameters.
// @Tags SIESA
// @Accept json
// @Produce application/json
// @Param siesa body RequestExcelCreate true "Request body for creating JSON file"
// @Security ApiKeyAuth
// @Failure 400 {object} shared.Response "Bad Request response"
// @Failure 422 {object} shared.Response "Unprocessable Entity response"
// @Failure 500 {object} shared.Response "Internal Server Error response"
// @Router /siesa/JSON [post]
func (h *Handler) CreateJSON(ctx *gin.Context) {
	_, doc, date, orders, ierr := h.initSIESA(ctx, "json")
	if ierr != nil {
		ctx.JSON(ierr.Code, shared.ErrorResponse(ierr.Message))
		return
	}

	// Call the HandleSIESAIntegration function to get the Excel file as a byte slice
	resp, err := h.service.HandleSIESAIntegrationJSON(doc, date, orders)
	if err != nil {
		shared.LogError("error handling SIESA integration", LogHandler, "Create", err, orders)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	defer h.service.UpdateDocumentStatus(doc, "success", "")

	ctx.JSON(http.StatusOK, shared.SuccessResponse(resp))
}

// Run handles a request to run the integration
// @Summary Run the integration
// @Description Run the integration
// @Tags SIESA
// @Accept json
// @Produce json
// @Param siesa body RequestExcelCreate true "Request body for running the integration"
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /siesa/run [post]
func (h *Handler) Run(ctx *gin.Context) {
	_, doc, date, orders, ierr := h.initSIESA(ctx, "integration")
	if ierr != nil {
		ctx.JSON(ierr.Code, shared.ErrorResponse(ierr.Message))
		return
	}

	resp, err := h.service.HandleSIESAIntegrationJSON(doc, date, orders)
	if err != nil {
		shared.LogError("error handling SIESA integration", LogHandler, "Create", err, orders)
		defer h.service.UpdateDocumentStatus(doc, "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	// since the integration can take a while, we run it in a goroutine and update
	// the document status when it finishes
	go func() {
		err := h.service.RunIntegration(resp)
		if err != nil {
			shared.LogError("error running integration", LogHandler, "Run", err, resp)
			h.service.UpdateDocumentStatus(doc, "error", err.Error())

			return
		}

		h.service.UpdateDocumentStatus(doc, "success", "")
	}()

	ctx.JSON(http.StatusOK, shared.SuccessResponse(doc))
}

// GetRunHistory handles a request to get the run history
// @Summary Get the run history
// @Description Get the run history
// @Tags SIESA
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param limit query number false "Limit of results"
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /siesa/history [get]
func (h *Handler) GetRunHistory(ctx *gin.Context) {
	sLimit := ctx.Query("limit")
	limit := 100
	if v, err := strconv.Atoi(sLimit); sLimit != "" && err == nil {
		limit = v
	}

	history, err := h.service.GetRunHistory(limit)
	if err != nil {
		shared.LogError("error getting run history", LogHandler, "GetRunHistory", err)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(history))
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
	req, doc, date, orders, ierr := h.initSIESA(ctx, "excel")
	if ierr != nil {
		ctx.JSON(ierr.Code, shared.ErrorResponse(ierr.Message))
		return
	}

	// Call the HandleSIESAIntegration function to get the Excel file as a byte slice
	excelFile, err := h.service.HandleSIESAIntegration(doc, date, orders)
	if err != nil {
		shared.LogError("error handling SIESA integration", LogHandler, "Create", err, orders)
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(ErrorInternalServer))
		return
	}

	filename := fmt.Sprintf("SIESA-%s-%s", strings.Join(req.LocationIDs, "-"), date.Format("2006-01-02"))

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")

	defer h.service.UpdateDocumentStatus(doc, "success", filename)

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

func (h *Handler) initSIESA(ctx *gin.Context, docType string) (*RequestExcelCreate, *SiesaDocument, time.Time, []PopappOrder, *shared.GinError) {
	var requestBody RequestExcelCreate
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)

		return nil, nil, time.Time{}, nil, &shared.GinError{Code: http.StatusBadRequest, Message: ErrorBadRequest}
	}

	if len(requestBody.LocationIDs) == 0 {
		return &requestBody, nil, time.Time{}, nil, &shared.GinError{Code: http.StatusBadRequest, Message: "location_ids is required"}
	}

	date, err := time.Parse("2006-01-02", requestBody.StartDate)
	if err != nil {
		shared.LogError("error parsing date", LogHandler, "Create", err, requestBody.StartDate)

		return &requestBody, nil, time.Time{}, nil, &shared.GinError{Code: http.StatusBadRequest, Message: "invalid date format for start_date " + err.Error()}
	}

	response, err := GetOrders(requestBody.StartDate, requestBody.EndDate, requestBody.LocationIDs)
	if err != nil {
		shared.LogError("error getting orders", LogHandler, "Create", err, response)
		return &requestBody, nil, time.Time{}, nil, &shared.GinError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	var orders []PopappOrder
	if err := json.Unmarshal([]byte(response), &orders); err != nil {
		shared.LogError("error unmarshalling orders", LogHandler, "Create", err, response)

		return &requestBody, nil, time.Time{}, nil, &shared.GinError{Code: http.StatusInternalServerError, Message: ErrorInternalServer}
	}

	if len(orders) == 0 {
		return &requestBody, nil, time.Time{}, nil, &shared.GinError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("No orders found at %s to %s in stores %+v", requestBody.StartDate, requestBody.EndDate, requestBody.LocationIDs),
		}
	}

	doc, err := h.service.GetDocument(requestBody.LocationIDs, requestBody.StartDate, requestBody.EndDate, len(orders), docType)
	if err != nil {
		shared.LogError("error creating document", LogHandler, "Create", err, doc)

		return &requestBody, nil, date, nil, &shared.GinError{Code: http.StatusInternalServerError, Message: ErrorInternalServer}
	}

	return &requestBody, doc, date, orders, nil
}
