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

// Log
type Log interface {
	Info(ctx context.Context, info domain.LogInfo)
	Error(ctx context.Context, info domain.LogError)
}

// Trace
type Trace interface {
	InitTracer(scopeName string) TraceUtil
	Tag() Tag
	Code() Code
}
type TraceUtil interface {
	ParaentSpan(ctx context.Context, name string) ParentSpanUtil
}

type ParentSpanUtil interface {
	End()
	GetSpan() GetSpan
	ChildSpan(name string) ParentSpanUtil
	SetTag(tag domain.TracerTag)
	AddEvent(eventName string, groupTag ...domain.TracerTag)
	SetStatus(status domain.TracerStatus, description string)
	SetName(name string)
	SetLink(spanID string, groupTag ...domain.TracerTag)
}

type Tag interface {
	String(key string, value string) domain.TracerTag
	Int64(key string, value int64) domain.TracerTag
}

type Code interface {
	OK() domain.TracerStatus
	Error() domain.TracerStatus
}

type GetSpan interface {
	ID() string
	Context() context.Context
}
