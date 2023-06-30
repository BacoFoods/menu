package taxes

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/tax/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create to handle a request for create taxes
// @Tags Tax
// @Summary To create a tax
// @Description To create a tax
// @Param tax body Tax true "tax request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Tax}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /tax [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody Tax
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	tax, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating tax", LogHandler, "Create", err, tax)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingTax))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(tax))
}

// Find to handle a request for find taxes
// @Tags Tax
// @Summary To find tax
// @Description To find tax
// @Param country_id query string false "country id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Tax}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /tax [get]
func (h *Handler) Find(ctx *gin.Context) {
	query := make(map[string]string)
	countryID := ctx.Query("country_id")
	if countryID != "" {
		query["country_id"] = countryID
	}
	tax, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error getting all tax", LogHandler, "Find", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingTax))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(tax))
}

// Get to handle a request for find tax by id
// @Tags Tax
// @Summary To find tax by id
// @Description To find tax by id
// @Param id path string true "tax id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Tax}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /tax/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	taxID := ctx.Param("id")
	tax, err := h.service.Get(taxID)
	if err != nil {
		shared.LogError("error getting tax", LogHandler, "Get", err, taxID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingTax))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(tax))
}

// Update to handle a request for update tax
// @Tags Tax
// @Summary To update tax
// @Description To update tax plan
// @Param tax body Tax true "tax to update"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Tax}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /tax [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var requestBody Tax
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogError("error getting tax request body", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	tax, err := h.service.Update(requestBody)
	if err != nil {
		shared.LogError("error updating tax", LogHandler, "Update", err, tax)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingTax))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(tax))
}

// Delete to handle a request for delete tax
// @Tags Tax
// @Summary To delete a tax
// @Description To delete a tax
// @Param id path string true "taxID"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Tax}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /tax/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	taxID := ctx.Param("id")
	tax, err := h.service.Delete(taxID)
	if err != nil {
		shared.LogError("error deleting tax", LogHandler, "Delete", err, taxID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingTax))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(tax))
}
