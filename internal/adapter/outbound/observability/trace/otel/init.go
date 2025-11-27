package otel

import (
	"Badminton-Hub/util"
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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
