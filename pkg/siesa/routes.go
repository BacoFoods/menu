package siesa

import (
	"github.com/BacoFoods/menu/pkg/shared"
)

type Routes struct {
	handler *Handler
}

func (r Routes) RegisterRoutes(private *shared.CustomRoutes) {
	// Get siesa file in json format
	private.POST("/siesa/json", r.handler.CreateJSON)

	// Get siesa file in a xlsx file
	private.POST("/siesa", r.handler.Create)

	// Run the integration
	private.POST("/siesa/run", r.handler.Run)

	// Get the run history
	private.GET("/siesa/history", r.handler.GetRunHistory)

	private.GET("/siesa/locales", r.handler.GetLocales)
	private.GET("/siesa/reference", r.handler.FindReferences)
	private.POST("/siesa/reference", r.handler.CreateReference)
	private.DELETE("/siesa/reference/:id", r.handler.DeleteReference)
	private.PATCH("/siesa/reference/:id", r.handler.UpdateReference)

}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}
