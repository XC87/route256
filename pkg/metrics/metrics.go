package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InitMetricsServer(mux *http.ServeMux) {
	if mux != nil {
		mux.Handle("/metrics", promhttp.Handler())
		return
	}

	http.Handle("/metrics", promhttp.Handler())
}
