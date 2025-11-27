package otel

import "go.opentelemetry.io/otel/trace"

func traceIDWithString(traceID string) (trace.TraceID, error) {
	return trace.TraceIDFromHex(traceID)
}

func spanIDWithString(spanID string) (trace.SpanID, error) {
	return trace.SpanIDFromHex(spanID)
}
