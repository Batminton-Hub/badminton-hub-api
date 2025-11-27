package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type Observability interface {
	Metrics() Metrics
	Log() Log
	Trace() Trace
}

// Metrics
type Metrics interface {
	GetMetrics(info domain.MetricsHttp)
	Counter(info domain.MetricsCounter) MetricsCounterUtil
	Gauge(info domain.MetricsGauge) MetricsGaugeUtil
}

type MetricsCounterUtil interface {
	Inc()
	Add(value float64) error
}

type MetricsGaugeUtil interface {
	Inc()
	Dec()
	Set(value float64)
	Add(value float64)
	Sub(value float64)
}

type MetricsHistogramUtil interface {
	Observe(value float64)
}

// Log
type Log interface {
	Info(ctx context.Context, info domain.LogInfo)
	Error(ctx context.Context, info domain.LogError)
}

// Trace
type Trace interface {
	SetScope(scopeName string) Span
	NewContext(traceID string, spanID string) (context.Context, error)
	Tag() Tag
	Code() Code
}

type Span interface {
	CreateSpan(ctx context.Context, name string) SpanUtil
	ConnectSpan(ctx context.Context) SpanUtil
}

type SpanUtil interface {
	End()
	AddSpan(name string) SpanUtil
	GetTraceID() string
	GetSpanID() string
	SetTag(tags ...domain.TracerTag)
	AddEvent(name string, tags ...domain.TracerTag)
	AddLink(ctx context.Context, tags ...domain.TracerTag)
	SetStatus(status domain.TracerStatus, description string)
}

type Tag interface {
	String(key string, value string) domain.TracerTag
	Int64(key string, value int64) domain.TracerTag
	Bool(key string, value bool) domain.TracerTag
}

type Code interface {
	OK() domain.TracerStatus
	Error() domain.TracerStatus
}
