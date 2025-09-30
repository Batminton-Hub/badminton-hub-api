package port

import "Badminton-Hub/internal/core/domain"

type Observability interface {
	Metrics() Metrics
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
