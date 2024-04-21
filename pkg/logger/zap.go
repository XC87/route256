package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func InitLogger(logLevelStr string, serviceName string) func() {
	zapConfig := zap.NewProductionConfig()
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(logLevelStr)); err != nil {
		log.Fatal(err)
	}
	lvl := zap.NewAtomicLevelAt(logLevel)

	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.Level = lvl

	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger = logger.With(zap.Field{Key: "service", Type: zapcore.StringType, String: serviceName})
	zap.ReplaceGlobals(logger)
	return func() {
		logger.Sync()
	}
}

func GetTraceFields(ctx context.Context) []zap.Field {
	spanCtx := trace.SpanContextFromContext(ctx)
	traceID := extractTraceID(spanCtx)
	spanID := extractSpanID(spanCtx)
	if traceID != "" && spanID != "" {
		return []zap.Field{zap.String("trace_id", traceID), zap.String("span_id", spanID)}
	}
	return []zap.Field{}
}

func extractSpanID(spanCtx trace.SpanContext) string {
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return traceID.String()
	}
	return ""
}

func extractTraceID(spanCtx trace.SpanContext) string {
	if spanCtx.HasSpanID() {
		spanID := spanCtx.SpanID()
		return spanID.String()
	}
	return ""
}
