package healthcheck

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler for health check handler
type Handler struct{}

// NewHandler create a new handler instance
func NewHandler() Handler {
	return Handler{}
}

// HealthCheck for handle health check request
func (h *Handler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}
