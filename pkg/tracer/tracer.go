package tracer

import (
	"context"
	"fmt"
	jaegerprop "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func InitTracer(tracerUrl string, serviceName string) (func(ctx context.Context), error) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(tracerUrl),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(1*time.Second)),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String("development"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(jaegerprop.Jaeger{})

	return func(ctx context.Context) {
		tracerProvider.ForceFlush(ctx)
		tracerProvider.Shutdown(ctx)
	}, nil
}

func HandleMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return HandlerMiddleware(next)
}

func HandlerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/health") || (r.URL.Path == "/metrics") {
			next.ServeHTTP(w, r)
			return
		}
		ctx, span := otel.Tracer("default").Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL))
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
