package account

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler = "pkg/account/handler.go"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create to handle a request to create an account
// @Tags Account
// @Summary To create an account
// @Description To create an account
// @Param account body Account true "account request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /account [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody RequestAccount
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	account, err := h.service.Create(requestBody.ToAccount())
	if err != nil {
		shared.LogError("error creating account", LogHandler, "Create", err, account)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountCreating))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(account))
}

// Login to handle a request to login
// @Tags Account
// @Summary To login
// @Description To login
// @Param account body Account true "account request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /account/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var request RequestLogin
	if err := ctx.BindJSON(&request); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Login", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	account, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		shared.LogError("error logging in", LogHandler, "Login", err, account)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountLogin))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(account))
}

// Delete to handle a request to delete an account
// @Tags Account
// @Summary To delete an account
// @Description To delete an account
// @Param id path string true "account id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /account/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		shared.LogWarn("warning binding request fail", LogHandler, "Delete", nil)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	if err := h.service.Delete(id); err != nil {
		shared.LogError("error deleting account", LogHandler, "Delete", err, id)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountDeleting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse("account deleted"))
}
