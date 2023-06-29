package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORSMiddleware for handle cors request from client
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}

// Authentication for handle authentication request from client
func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.Request.Header.Get("Authorization") // TODO: improve security access
		ctx.Next()
	}
}
