package router

import (
	"encoding/base64"
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"

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

// AuthMiddleware for handle authentication request from client
func AuthMiddleware(validator *idtoken.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		shared.LogInfo("tokenString", LogMiddleware, "Authentication", nil, tokenString)

		if tokenString == "" {
			shared.LogWarn("token is empty", LogMiddleware, "Authentication", nil)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if GoogleAuth(tokenString, ctx, validator) || CerebroAuth(tokenString, ctx) {
			ctx.Next()
			return
		}

		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func GoogleAuth(credential string, ctx *gin.Context, validator *idtoken.Validator) bool {
	_, err := validator.Validate(ctx, credential, "")
	if err != nil {
		shared.LogWarn("failed to validate google token", LogMiddleware, "GoogleAuth", err, credential)
		return false
	}

	return true
}

func CerebroAuth(credential string, ctx *gin.Context) bool {
	secretKey, err := base64.StdEncoding.DecodeString(internal.Config.TokenSecret)
	if err != nil {
		shared.LogError("error decoding jwt key", LogMiddleware, "Authentication", err, internal.Config.TokenSecret)
		return false
	}

	token, err := jwt.Parse(credential, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			shared.LogWarn("method of signature is invalid", LogMiddleware, "Authentication", nil)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})
	if err != nil {
		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return false
	}

	ctx.Set("account_uuid", token.Claims.(jwt.MapClaims)["uuid"])
	ctx.Set("account_role", token.Claims.(jwt.MapClaims)["role"])
	ctx.Set("account_name", token.Claims.(jwt.MapClaims)["name"])
	ctx.Set("brand_id", token.Claims.(jwt.MapClaims)["brand"])
	ctx.Set("brand_name", token.Claims.(jwt.MapClaims)["brand_name"])
	ctx.Set("channel_id", token.Claims.(jwt.MapClaims)["channel"])
	ctx.Set("channel_name", token.Claims.(jwt.MapClaims)["channel_name"])
	ctx.Set("store_id", token.Claims.(jwt.MapClaims)["store"])
	ctx.Set("store_name", token.Claims.(jwt.MapClaims)["store_name"])

	return true
}
