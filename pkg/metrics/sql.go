package metrics

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
	"time"
)

type MetricQueryInfo struct {
	SelectCounter prometheus.Counter
	UpdateCounter prometheus.Counter
	DeleteCounter prometheus.Counter
	SelectTime    *prometheus.HistogramVec
	UpdateTime    *prometheus.HistogramVec
	DeleteTime    *prometheus.HistogramVec
}

type SqlMetric struct {
	metrics *MetricQueryInfo
}

type Timer struct {
	begin    time.Time
	observer *prometheus.HistogramVec
}

func NewSqlMetrics() *SqlMetric {
	return &SqlMetric{
		metrics: createDbQueryMetrics(),
	}
}

func createDbQueryMetrics() *MetricQueryInfo {
	m := &MetricQueryInfo{}
	m.SelectCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "select_query_total",
		Help: "Total number of select queries",
	})
	m.UpdateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "update_query_total",
		Help: "Total number of update queries",
	})
	m.DeleteCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delete_query_total",
		Help: "Total number of delete queries",
	})
	m.SelectTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "select_duration_seconds",
		Help:    "Duration of select queries in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})
	m.UpdateTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "update_duration_seconds",
		Help:    "Duration of update queries in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})
	m.DeleteTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "delete_duration_seconds",
		Help:    "Duration of delete queries in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})
	return m
}

func NewTimer(o *prometheus.HistogramVec) *Timer {
	return &Timer{
		time.Now(),
		o,
	}
}

func (tracer *SqlMetric) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	var queryTimer *Timer

	if strings.Contains(data.SQL, "select") {
		tracer.metrics.SelectCounter.Inc()
		queryTimer = NewTimer(tracer.metrics.SelectTime)
	} else if strings.Contains(data.SQL, "update") {
		tracer.metrics.UpdateCounter.Inc()
		queryTimer = NewTimer(tracer.metrics.UpdateTime)
	} else if strings.Contains(data.SQL, "delete") {
		tracer.metrics.DeleteCounter.Inc()
		queryTimer = NewTimer(tracer.metrics.DeleteTime)
	}

	return context.WithValue(ctx, "queryTimer", queryTimer)
}

func (tracer *SqlMetric) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	queryTimer, ok := ctx.Value("queryTimer").(*Timer)
	if !ok || queryTimer == nil {
		return
	}
	status := "success"
	if data.Err != nil {
		status = "error"
	}
	queryTimer.observer.WithLabelValues(status).Observe(time.Since(queryTimer.begin).Seconds())

}
