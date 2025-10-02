package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TraceUtil struct {
	Tracer trace.Tracer
}

type SpanUtil struct {
	Tracer     trace.Tracer
	Ctx        context.Context
	Span       trace.Span
	Attributes []attribute.KeyValue
}

type Tag struct{}

type Code struct{}

type Span struct {
	Span trace.Span
}
