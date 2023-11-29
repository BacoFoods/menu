package scheduler

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.GET("/schedules", r.handler.Find)
	private.POST("/schedules", r.handler.Create)
	private.PATCH("/schedules", r.handler.Update)
	private.DELETE("/schedules", r.handler.Delete)
	private.GET("/schedules/store/:id/today", r.handler.TodayStore)
	private.GET("/schedules/brand/:id/today", r.handler.TodayBrand)
	private.POST("/schedules/store/:id/enable", r.handler.EnableStore)
	private.POST("/schedules/holiday", r.handler.Holiday)
	private.PATCH("/schedules/holiday", r.handler.UpdateHoliday)
	private.DELETE("/schedules/holiday", r.handler.DeleteHoliday)
	private.GET("/schedules/holiday", r.handler.FindHoliday)
}
