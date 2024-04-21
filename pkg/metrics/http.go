package metrics

import (
	status_writer "github.com/joegasewicz/status-writer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MetricInfo struct {
	RequestCounter *prometheus.CounterVec
	ResponseTime   prometheus.Histogram
}

func CreateMetricMiddleware(serviceName string, handlerName string, handleURL string) func(next http.HandlerFunc) http.HandlerFunc {
	metrics := CreateRequestMetrics(serviceName, handlerName, handleURL)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func(start time.Time) {
				metrics.ResponseTime.Observe(time.Since(start).Seconds())
			}(time.Now())
			sw := status_writer.New(w)
			next(sw, r)
			statusCode := sw.Status

			metrics.RequestCounter.WithLabelValues(handleURL, strconv.Itoa(statusCode)).Inc()
		}
	}
}

func CreateRequestMetrics(serviceName string, handlerName string, handleURL string) MetricInfo {
	handlerName = strings.ToLower(handlerName)
	counterName := handlerName + "_requests_total"
	requestCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "cart",
			Name:      counterName,
			Help:      "Total amount of requests made to " + handleURL + " handler. Example: rate(" + counterName + "[1m])",
		},
		[]string{"url", "status"},
	)
	responseName := handlerName + "_response_time"
	responseTime := promauto.NewHistogram(prometheus.HistogramOpts{
		Subsystem: serviceName,
		Name:      responseName,
		Buckets:   prometheus.DefBuckets,
	})

	return MetricInfo{
		RequestCounter: requestCounter,
		ResponseTime:   responseTime,
	}
}
