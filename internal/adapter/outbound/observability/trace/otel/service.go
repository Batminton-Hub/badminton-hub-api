package otel

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TraceUtil
func (o *Otel) SetScope(scopeName string) port.Span {
	return &Span{
		ScopeName: scopeName,
	}
}

func (o *Otel) NewContext(traceIDStr string, spanIDStr string) (context.Context, error) {
	traceID, err := traceIDWithString(traceIDStr)
	if err != nil {
		return nil, err
	}
	spanID, err := spanIDWithString(spanIDStr)
	if err != nil {
		return nil, err
	}
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	newCtx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	return newCtx, nil
}

func (o *Otel) Tag() port.Tag {
	return &Tag{}
}

func (o *Otel) Code() port.Code {
	return &Code{}
}

// Span
func (o *Span) CreateSpan(ctx context.Context, name string) port.SpanUtil {
	tracer := otel.Tracer(o.ScopeName)
	newCtx, span := tracer.Start(ctx, name)
	return &SpanUtil{
		Context: newCtx,
		Tracer:  tracer,
		Span:    span,
	}
}

func (o *Span) ConnectSpan(ctx context.Context) port.SpanUtil {
	tracer := otel.Tracer(o.ScopeName)
	span := trace.SpanFromContext(ctx)
	return &SpanUtil{
		Context: ctx,
		Tracer:  tracer,
		Span:    span,
	}
}

// SpanUtil
func (s *SpanUtil) End() {
	s.Span.End()
}

func (s *SpanUtil) AddSpan(name string) port.SpanUtil {
	newCtx, span := s.Tracer.Start(s.Context, name)
	return &SpanUtil{
		Context: newCtx,
		Tracer:  s.Tracer,
		Span:    span,
	}
}

func (s *SpanUtil) GetTraceID() string {
	return s.Span.SpanContext().TraceID().String()
}

func (s *SpanUtil) GetSpanID() string {
	return s.Span.SpanContext().SpanID().String()
}

func (s *SpanUtil) SetTag(tags ...domain.TracerTag) {
	buildTag := []attribute.KeyValue{}
	for _, tag := range tags {
		buildTag = append(buildTag, tag.Attribute.(attribute.KeyValue))
	}
	if len(buildTag) > 0 {
		s.Span.SetAttributes(buildTag...)
	}
}

func (s *SpanUtil) AddEvent(name string, tags ...domain.TracerTag) {
	buildTag := []attribute.KeyValue{}
	for _, tag := range tags {
		buildTag = append(buildTag, tag.Attribute.(attribute.KeyValue))
	}
	if len(buildTag) > 0 {
		s.Span.AddEvent(name, trace.WithAttributes(buildTag...))
	}
}

func (s *SpanUtil) AddLink(ctx context.Context, tags ...domain.TracerTag) {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	buildTag := []attribute.KeyValue{}
	for _, tag := range tags {
		buildTag = append(buildTag, tag.Attribute.(attribute.KeyValue))
	}
	link := trace.Link{
		SpanContext: spanContext,
		Attributes:  buildTag,
	}
	if len(buildTag) > 0 {
		s.Span.AddLink(link)
	}
}

func (s *SpanUtil) SetStatus(status domain.TracerStatus, description string) {
	s.Span.SetStatus(status.Code.(codes.Code), description)
}

func (s *SpanUtil) SetName(name string) {
	s.Span.SetName(name)
}

// TagUtil
func (t *Tag) String(key string, value string) domain.TracerTag {
	return domain.TracerTag{
		TypeVal:   domain.String,
		Attribute: attribute.String(key, value),
	}
}
func (t *Tag) Int64(key string, value int64) domain.TracerTag {
	return domain.TracerTag{
		TypeVal:   domain.Int64,
		Attribute: attribute.Int64(key, value),
	}
}
func (t *Tag) Bool(key string, value bool) domain.TracerTag {
	return domain.TracerTag{
		TypeVal:   domain.Bool,
		Attribute: attribute.Bool(key, value),
	}
}

// CodeUtil
func (c *Code) OK() domain.TracerStatus {
	return domain.TracerStatus{
		Code: codes.Ok,
	}
}
func (c *Code) Error() domain.TracerStatus {
	return domain.TracerStatus{
		Code: codes.Error,
	}
}
