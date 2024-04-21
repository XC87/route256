package tracer

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type SqlTracerGroup struct {
	Tracer []SqlTracer
}
type SqlMetric struct {
	metrics SqlTracer
}

type SqlTracer interface {
	TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context
	TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData)
}

func NewSqlTracerGroup(tracer ...SqlTracer) *SqlTracerGroup {
	return &SqlTracerGroup{
		Tracer: tracer,
	}
}

func (tracerGroup *SqlTracerGroup) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, sqlTracer := range tracerGroup.Tracer {
		ctx = sqlTracer.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

func (tracerGroup *SqlTracerGroup) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, sqlTracer := range tracerGroup.Tracer {
		sqlTracer.TraceQueryEnd(ctx, conn, data)
	}
}
