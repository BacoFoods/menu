package router

import (
	"encoding/base64"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogMiddleware string = "pkg/router/middleware"

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
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			shared.LogWarn("token is empty", LogMiddleware, "Authentication", nil)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		secretKey, err := base64.StdEncoding.DecodeString(internal.Config.TokenSecret)
		if err != nil {
			shared.LogError("error decoding jwt key", LogMiddleware, "Authentication", err, internal.Config.TokenSecret)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				shared.LogWarn("method of signature is invalid", LogMiddleware, "Authentication", nil)
				ctx.AbortWithStatus(http.StatusUnauthorized)
			}

			return secretKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Next()
	}
}
