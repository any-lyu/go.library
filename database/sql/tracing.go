package sql

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracinglog "github.com/opentracing/opentracing-go/log"
)

// startTracingSpan start a span from context before any db operation
func startTracingSpan(ctx context.Context, operationName, query string) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		// log.ErrorContext(ctx, "cannot found tracing span from context") // TODO: enable this
		return nil
	}
	span := opentracing.GlobalTracer().StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	ext.SpanKindRPCClient.Set(span)
	ext.DBType.Set(span, "database")
	ext.DBStatement.Set(span, query)
	return span
}

// finishTracingSpan finish a span after db operation
func finishTracingSpan(span opentracing.Span, err error) {
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracinglog.String("event", "error"), tracinglog.String("message", err.Error()))
	}
	span.Finish()
}
