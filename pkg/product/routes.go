package product

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *shared.CustomRoutes) {
	// Products
	private.GET("/product", r.handler.Find)
	private.GET("/product/:id", r.handler.Get)
	private.POST("/product", r.handler.Create)
	private.PATCH("/product/:id", r.handler.Update)
	private.DELETE("/product/:id", r.handler.Delete)
	private.POST("/product/:id/modifier/:modifierID", r.handler.AddModifier)
	private.DELETE("/product/:id/modifier/:modifierID", r.handler.RemoveModifier)
	private.GET("/product/:id/overrider", r.handler.GetOverridersByField)
	private.PATCH("/product/:id/overrider/update-all", r.handler.UpdateAllOverriders)
	private.GET("/product/:id/category", r.handler.GetCategories)

	// Modifiers
	private.GET("/modifier", r.handler.ModifierFind)
	private.POST("/modifier", r.handler.ModifierCreate)
	private.POST("/modifier/:id/product/:productID", r.handler.ModifierAddProduct)
	private.PATCH("/modifier/:id", r.handler.ModifierUpdate)
	private.DELETE("/modifier/:id/product/:productID", r.handler.ModifierRemoveProduct)

	// Overriders
	private.GET("/overrider", r.handler.OverriderFind)
	private.GET("/overrider/:id", r.handler.OverriderGet)
	private.POST("/overrider", r.handler.OverriderCreate)
	private.PATCH("/overrider/:id", r.handler.OverriderUpdate)
	private.DELETE("/overrider/:id", r.handler.OverriderDelete)
}
