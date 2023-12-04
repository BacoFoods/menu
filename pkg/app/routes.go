package app

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{
		handler: handler,
	}
}

func (r *Routes) RegisterRoutes(routes *gin.RouterGroup) {
	routes.GET("/download/windows/:version", r.handler.DownloadWindows)
	routes.GET("/download/apk/:version", r.handler.DownloadAPK)
}
