package core_util

import "Badminton-Hub/internal/core/port"

type Observability struct {
	metrics port.Metrics
}

func NewObservability(
	metrics port.Metrics,
) *Observability {
	return &Observability{
		metrics: metrics,
	}
}

func (o *Observability) Metrics() port.Metrics {
	return o.metrics
}
