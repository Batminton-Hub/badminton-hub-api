package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func traceID(ctx context.Context) trace.TraceID {
	var traceID trace.TraceID
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID = span.SpanContext().TraceID()
	}
	return traceID
}

func spanContext(traceID trace.TraceID) trace.SpanContext {
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  trace.SpanID{},
	})
	return spanCtx
}

func stringToTraceID(traceID string) (trace.TraceID, error) {
	newTraceID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return trace.TraceID{}, err
	}
	return newTraceID, nil
}
