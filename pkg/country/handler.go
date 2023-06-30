package country

import (
	"net/http"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler string = "pkg/country/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create to handle a request to create a country
// @Tags Country
// @Summary To create a country
// @Description To create a country
// @Param country body Country true "country request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Country}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /country [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody Country
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	country, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating country", LogHandler, "Create", err, country)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingCountry))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(country))
}

// Find to handle a request to find all country
// @Tags Country
// @Summary To find country
// @Description To find country
// @Param code query string false "country code"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Country}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /country [get]
func (h *Handler) Find(ctx *gin.Context) {
	query := make(map[string]string)
	code := ctx.Query("code")
	if code != "" {
		query["iso_code"] = ctx.Query("code")
	}
	country, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error getting all country", LogHandler, "Find", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingCountry))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(country))
}

// Get to handle a request to find country by id
// @Tags Country
// @Summary To find country by id
// @Description To find country by id
// @Param id path string true "country id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Country}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /country/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	countryID := ctx.Param("id")
	country, err := h.service.Get(countryID)
	if err != nil {
		shared.LogError("error getting country", LogHandler, "Get", err, countryID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingCountry))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(country))
}

// Update to handle a request to update country
// @Tags Country
// @Summary To update country
// @Description To update country plan
// @Param country body Country true "country to update"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Country}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /country [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var requestBody Country
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogError("error getting country request body", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	country, err := h.service.Update(requestBody)
	if err != nil {
		shared.LogError("error updating country", LogHandler, "Update", err, country)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingCountry))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(country))
}

// Delete to handle a request to delete a country
// @Tags Country
// @Summary To delete a country
// @Description To delete a country
// @Param id path string true "countryID"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Country}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /country/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	countryID := ctx.Param("id")
	country, err := h.service.Delete(countryID)
	if err != nil {
		shared.LogError("error deleting country", LogHandler, "Delete", err, countryID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingCountry))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(country))
}
