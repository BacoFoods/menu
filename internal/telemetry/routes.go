package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Routes struct {
}

type DataPoint struct {
	Event string         `json:"event"`
	Value int64          `json:"value"`
	Data  map[string]any `json:"data"`
}

func (r *Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/telemetry", postTelemetry)
}

func postTelemetry(c *gin.Context) {
	var data []DataPoint
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// we early return OK
	c.JSON(http.StatusOK, gin.H{})

	for _, dt := range data {
		tags := make(map[string]string)
		for k, v := range dt.Data {
			tags[k] = fmt.Sprintf("%v", v)
		}

		point := TelemetryPoint{
			Tags:        tags,
			Measurement: float64(dt.Value),
			End:         time.Now(),
		}

		go report(dt.Event, &point)
	}
}
