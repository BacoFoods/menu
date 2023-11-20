package siesa

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/siesa", r.handler.Create)
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}
