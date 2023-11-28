package telemetry

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/BacoFoods/menu/internal"
	"github.com/sirupsen/logrus"
)

const (
	TelemetryServiceName     = "menu"
	TelemetryMeasurementName = "response_time"
)

var (
	writeUrl string
	headers  map[string]string
	client   http.Client
)

func init() {
	headers = map[string]string{
		"Authorization": fmt.Sprintf("Token %s", internal.Config.InfluxToken),
	}

	host := fmt.Sprintf("http://%s:%s", internal.Config.InfluxHost, internal.Config.InfluxPort)
	writeUrl = fmt.Sprintf("%s/write?db=%s", host, internal.Config.InfluxDB)
	client = http.Client{
		Timeout: 1 * time.Second,
	}
}

func post(payload string) {
	req, err := http.NewRequest("POST", writeUrl, strings.NewReader(payload))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		logrus.Warn("Error sending telemetry", err)
		return
	}

	if res != nil {
		defer res.Body.Close()
		_, _ = io.Copy(io.Discard, res.Body)
	}
}

func report(metric string, point *TelemetryPoint) {
	if point == nil {
		return
	}

	data := []string{metric}
	for k, v := range point.Tags {
		data = append(data, fmt.Sprintf("%s=%s", k, v))
	}

	ts := point.End.UnixNano()

	payload := fmt.Sprintf("%s value=%f %d", strings.Join(data, ","), point.Measurement, ts)

	post(payload)
}

type TelemetryPoint struct {
	Tags        map[string]string
	Start       time.Time
	End         time.Time
	Measurement float64
}

func (t *TelemetryPoint) Done(status int, m int64) {
	if t == nil {
		return
	}

	if t.Tags == nil {
		t.Tags = make(map[string]string)
	}

	t.Tags["response_status"] = fmt.Sprintf("%d", status)

	t.End = time.Now()
	t.Measurement = float64(m)

	report(TelemetryMeasurementName, t)
}

func StartResponse(method, url string, now time.Time) *TelemetryPoint {
	return &TelemetryPoint{
		Tags: map[string]string{
			"method":  method,
			"url":     url,
			"service": TelemetryServiceName,
		},
		Start: now,
	}
}
