package discount

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/discount/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create to handle a request for create a discount
// @Tags Discount
// @Summary To create a discount
// @Description To create a discount
// @Param discount body Discount true "discount request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Discount}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /discount [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody Discount
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	discount, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating discount", LogHandler, "Create", err, discount)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingDiscount))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(discount))
}

// Find to handle a request for find all discount
// @Tags Discount
// @Summary To find discount
// @Description To find discount
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Discount}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /discount [get]
func (h *Handler) Find(ctx *gin.Context) {
	query := make(map[string]string)

	name := ctx.Query("name")
	if name != "" {
		query["name"] = ctx.Query("name")
	}

	discounts, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding discounts", LogHandler, "Find", err, discounts)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingDiscount))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(discounts))
}

// Get to handle a request for find a discount by id
// @Tags Discount
// @Summary To find discount by id
// @Description To find discount by id
// @Param id path string true "discount id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Discount}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /discount/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	discountID := ctx.Param("id")
	discount, err := h.service.Get(discountID)
	if err != nil {
		shared.LogError("error getting discount", LogHandler, "Get", err, discountID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingDiscount))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(discount))
}

// Update to handle a request for update a discount
// @Tags Discount
// @Summary To update discount
// @Description To update discount plan
// @Param discount body Discount true "discount to update"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Discount}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /discount [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var requestBody Discount
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogError("error getting discount request body", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	discount, err := h.service.Update(requestBody)
	if err != nil {
		shared.LogError("error updating discount", LogHandler, "Update", err, discount)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingDiscount))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(discount))
}

// Delete to handle a request for delete a discount
// @Tags Discount
// @Summary To delete a discount
// @Description To delete a discount
// @Param id path string true "discountID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Discount}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /discount/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	discountID := ctx.Param("id")
	discount, err := h.service.Delete(discountID)
	if err != nil {
		shared.LogError("error deleting discount", LogHandler, "Delete", err, discountID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingDiscount))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(discount))
}
