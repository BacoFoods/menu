package invoice

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *shared.CustomRoutes) {
	router.GET("/invoice/:id", r.handler.Get)
	router.GET("/invoice", r.handler.Find)
	router.GET("/invoice/:id/print", r.handler.Print)
	router.POST("/invoice/:id", r.handler.UpdateTip)
	router.POST("/invoice/:id/client/:clientID/add", r.handler.AddClient)
	router.POST("/invoice/:id/client/:clientID/remove", r.handler.RemoveClient)
	router.POST("/invoice/:id/split", r.handler.Split)

	router.GET("/discount-applied", r.handler.FindDiscountApplied)
	router.DELETE("/discount-applied/:id", r.handler.RemoveDiscountApplied)

	// DIAN Resolutions
	router.GET("/resolution", r.handler.FindResolution)
	router.POST("/resolution", r.handler.CreateResolution)
	router.PATCH("/resolution/:id", r.handler.UpdateResolution)
	router.DELETE("/resolution/:id", r.handler.DeleteResolution)
}
