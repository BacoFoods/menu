package cashaudit

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/cash-audit/orders-closed", r.handler.OrdersClosedValidation)
	router.POST("/cash-audit", r.handler.Create)
	router.GET("/cash-audit", r.handler.Get)
	router.POST("/cash-audit/confirm", r.handler.Confirm)
}
