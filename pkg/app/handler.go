package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const (
	zipContentType = "application/zip"
	apkContentType = "application/vnd.android.package-archive"

	logHandler = "pkg/app/handler.go"
)

type Service interface {
	DownloadWindows(version string) (string, io.ReadCloser, int64, error)
	DownloadAPK(version string) (string, io.ReadCloser, int64, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) DownloadWindows(c *gin.Context) {
	version := c.Param("version")

	filename, data, size, err := h.service.DownloadWindows(version)
	if err != nil {
		shared.LogError("error downloading", logHandler, "Download", err)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.DataFromReader(http.StatusOK, size, zipContentType, data, map[string]string{
		"Content-Disposition":           fmt.Sprintf("attachment; filename=%s", filename),
		"Content-Description":           "File Transfer",
		"Content-Type":                  zipContentType,
		"Content-Transfer-Encoding":     "binary",
		"Cache-Control":                 "must-revalidate",
		"Pragma":                        "public",
		"Access-Control-Expose-Headers": "Content-Disposition",
	})
}

func (h *Handler) DownloadAPK(c *gin.Context) {
	version := c.Param("version")

	filename, data, size, err := h.service.DownloadAPK(version)
	if err != nil {
		shared.LogError("error downloading", logHandler, "Download", err, data)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.DataFromReader(http.StatusOK, size, apkContentType, data, map[string]string{
		"Content-Disposition":           fmt.Sprintf("attachment; filename=%s", filename),
		"Content-Description":           "File Transfer",
		"Content-Type":                  apkContentType,
		"Content-Transfer-Encoding":     "binary",
		"Cache-Control":                 "must-revalidate",
		"Pragma":                        "public",
		"Access-Control-Expose-Headers": "Content-Disposition",
	})
}
