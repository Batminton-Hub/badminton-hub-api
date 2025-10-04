package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	ScopeName string
}

type SpanUtil struct {
	Context context.Context
	Tracer  trace.Tracer
	Span    trace.Span
}

type Tag struct{}

type Code struct{}
