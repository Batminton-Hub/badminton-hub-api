package core_util

import "Badminton-Hub/internal/core/port"

type Observability struct {
	metrics port.Metrics
	log     port.Log
}

func NewObservability(
	metrics port.Metrics,
	log port.Log,
) *Observability {
	return &Observability{
		metrics: metrics,
		log:     log,
	}
}

func (o *Observability) Metrics() port.Metrics {
	return o.metrics
}

func (o *Observability) Log() port.Log {
	return o.log
}
