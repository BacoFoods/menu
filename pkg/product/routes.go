package product

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	// Products
	router.GET("/product", r.handler.Find)
	router.GET("/product/:id", r.handler.Get)
	router.POST("/product", r.handler.Create)
	router.PATCH("/product/:id", r.handler.Update)
	router.DELETE("/product/:id", r.handler.Delete)
	router.POST("/product/:id/modifier/:modifierID", r.handler.AddModifier)
	router.DELETE("/product/:id/modifier/:modifierID", r.handler.RemoveModifier)

	// Modifiers
	router.GET("/modifier", r.handler.ModifierFind)
	router.POST("/modifier", r.handler.ModifierCreate)
	router.POST("/modifier/:id/product/:productID", r.handler.ModifierAddProduct)
	router.DELETE("/modifier/:id/product/:productID", r.handler.ModifierRemoveProduct)
}
