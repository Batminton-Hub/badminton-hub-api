package otel

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Otel struct{}

func NewOtel(serviceName string) *Otel {
	fmt.Println("Starting tracing")
	config := util.LoadConfig()

	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			// otlptracehttp.WithEndpoint("jaeger:4318"),
			// otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithEndpoint(config.TracerServerURL),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatalf("creating new exporter: %v", err)
		return nil
	}

	tracerprovider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)
	return &Otel{}
}

// TraceUtil
func (o *Otel) InitTracer(scopeName string) port.TraceUtil {
	fmt.Println("Starting tracer", scopeName)
	tracer := otel.Tracer(scopeName)
	return &TraceUtil{
		Tracer: tracer,
	}
}

func (o *Otel) Tag() port.Tag {
	return &Tag{}
}

func (o *Otel) Code() port.Code {
	return &Code{}
}

func (t *TraceUtil) ParaentSpan(ctx context.Context, name string) port.ParentSpanUtil {
	traceID := traceID(ctx)
	spanCtx := spanContext(traceID)
	newCtx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	newCtx, span := t.Tracer.Start(newCtx, name)
	fmt.Println("Starting span", span.SpanContext().SpanID().String())
	return SpanUtil{
		Tracer: t.Tracer,
		Ctx:    newCtx,
		Span:   span,
	}

}

// SpanUtil
func (s SpanUtil) End() {
	fmt.Println("Ending span", s.Span.SpanContext().SpanID().String())
	s.Span.End()
}

func (s SpanUtil) ChildSpan(spanName string) port.ParentSpanUtil {
	spanCtx := s.Span.SpanContext()
	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	ctx, childSpan := s.Tracer.Start(ctx, spanName)
	fmt.Println("Starting child span", childSpan.SpanContext().SpanID().String())
	return SpanUtil{
		Tracer: s.Tracer,
		Ctx:    ctx,
		Span:   childSpan,
	}
}

func (s SpanUtil) SetTag(tag domain.TracerTag) {
	s.Span.SetAttributes(tag.Attribute.(attribute.KeyValue))
}

func (s SpanUtil) AddEvent(eventName string, groupTag ...domain.TracerTag) {
	tags := []attribute.KeyValue{}
	for _, t := range groupTag {
		tags = append(tags, t.Attribute.(attribute.KeyValue))
	}
	eventOption := trace.WithAttributes(tags...)
	s.Span.AddEvent(eventName, eventOption)
}

func (s SpanUtil) SetStatus(status domain.TracerStatus, description string) {
	code := status.Code.(codes.Code)
	s.Span.SetStatus(code, description)
}

func (s SpanUtil) SetName(name string) {
	s.Span.SetName(name)
}

func (s SpanUtil) GetSpan() port.GetSpan {
	return &Span{
		Span: s.Span,
	}
}

func (s SpanUtil) SetLink(spanID string, groupTag ...domain.TracerTag) {
	tags := []attribute.KeyValue{}
	for _, t := range groupTag {
		tags = append(tags, t.Attribute.(attribute.KeyValue))
	}
	traceID, _ := stringToTraceID(spanID)
	link := trace.Link{
		SpanContext: spanContext(traceID),
		Attributes:  tags,
	}
	s.Span.AddLink(link)
}

// Tag
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

// Code
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

// GetSpan
func (s *Span) ID() string {
	return s.Span.SpanContext().SpanID().String()
}

func (s *Span) Context() context.Context {
	ctx := trace.ContextWithSpanContext(context.Background(), s.Span.SpanContext())
	return ctx
}
