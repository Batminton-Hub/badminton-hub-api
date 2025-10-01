package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type Observability interface {
	Metrics() Metrics
	Log() Log
}

// Metrics
type Metrics interface {
	GetMetrics(info domain.MetricsHttp)
	Counter(info domain.MetricsCounter) MetricsCounterUtil
	Gauge(info domain.MetricsGauge) MetricsGaugeUtil
}

// Log
type Log interface {
	Info(ctx context.Context, info domain.LogInfo)
	Error(ctx context.Context, info domain.LogError)
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
