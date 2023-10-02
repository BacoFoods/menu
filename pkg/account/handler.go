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
// @Param account body RequestAccount true "account request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /public/account [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody RequestAccount
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	account, err := h.service.Create(requestBody.ToAccount())
	if err != nil {
		shared.LogError("error creating account", LogHandler, "Create", err, requestBody)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountCreation))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(account))
}

// CreatePinUser to handle a request to create a pin user
// @Tags Account
// @Summary To create a pin user
// @Description To create a pin user
// @Param account body RequestPinUser true "account request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /account [post]
func (h *Handler) CreatePinUser(ctx *gin.Context) {
	var request RequestPinUser
	if err := ctx.BindJSON(&request); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "CreatePinUser", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	if request.Pin < 1000 || request.Pin > 9999 {
		shared.LogWarn("warning wrong pin length", LogHandler, "CreatePinUser", nil)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorAccountPinBadRequest))
		return
	}

	account, err := h.service.CreatePinUser(request.ToAccount())
	if err != nil {
		shared.LogError("error creating account", LogHandler, "CreatePinUser", err, request)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountCreation))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(account))
}

// Login to handle a request to login
// @Tags Account
// @Summary To login
// @Description To login
// @Param account body RequestLogin true "account request"
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 400 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /public/account/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var request RequestLogin
	if err := ctx.BindJSON(&request); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Login", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	account, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		shared.LogError("error logging in", LogHandler, "Login", err, request)
		ctx.JSON(http.StatusForbidden, shared.ErrorResponse(ErrorAccountLogin))
		return
	}

	jwt, err := account.JWT()
	if err != nil {
		shared.LogError("error generating jwt", LogHandler, "Login", err, account)
		ctx.JSON(http.StatusForbidden, shared.ErrorResponse(ErrorAccountLogin))
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}

// LoginPin to handle a request to login
// @Tags Account
// @Summary To login
// @Description To login
// @Param account body RequestLoginPin true "account request"
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 400 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /public/account/login/pin [post]
func (h *Handler) LoginPin(ctx *gin.Context) {
	var request RequestLoginPin
	if err := ctx.BindJSON(&request); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "LoginPin", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	if request.Pin < 1000 || request.Pin > 9999 {
		shared.LogWarn("warning wrong pin length", LogHandler, "LoginPin", nil)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorAccountLogin))
		return
	}

	account, err := h.service.LoginPin(request.Pin)
	if err != nil {
		shared.LogError("error logging in", LogHandler, "LoginPin", err, request)
		ctx.JSON(http.StatusForbidden, shared.ErrorResponse(ErrorAccountLogin))
		return
	}

	jwt, err := account.JWT()
	if err != nil {
		shared.LogError("error generating jwt", LogHandler, "LoginPin", err, account)
		ctx.JSON(http.StatusForbidden, shared.ErrorResponse(ErrorAccountLogin))
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}

// Delete to handle a request to delete an account
// @Tags Account
// @Summary To delete an account
// @Description To delete an account
// @Param id path string true "account id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
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

// Find to handle a request to find an account
// @Tags Account
// @Summary To find an account
// @Description To find an account
// @Param username query string false "account username"
// @Param email query string false "account email"
// @Param storeID query string false "account store id"
// @Param role query string false "account role"
// @Param brandID query string false "account brand id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /account [get]
func (h *Handler) Find(ctx *gin.Context) {
	filter := make(map[string]any)

	username := ctx.Query("username")
	if username != "" {
		filter["username"] = username
	}

	email := ctx.Query("email")
	if email != "" {
		filter["email"] = email
	}

	storeID := ctx.Query("storeID")
	if storeID != "" {
		filter["storeID"] = storeID
	}

	role := ctx.Query("role")
	if role != "" {
		filter["role"] = role
	}

	brandID := ctx.Query("brandID")
	if brandID != "" {
		filter["brandID"] = brandID
	}

	accounts, err := h.service.Find(filter)
	if err != nil {
		shared.LogError("error finding account", LogHandler, "Find", err, filter)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountFinding))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(accounts))
}

// Update to handle a request to update an account
// @Tags Account
// @Summary To update an account
// @Description To update an account
// @Param account body RequestAccountUpdate true "account request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Account}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /account [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var requestBody RequestAccountUpdate
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	account, err := h.service.Update(requestBody.ToAccount())
	if err != nil {
		shared.LogError("error updating account", LogHandler, "Update", err, requestBody)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAccountUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(account))
}
