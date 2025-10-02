package core_util

import "Badminton-Hub/internal/core/port"

type Observability struct {
	metrics port.Metrics
	log     port.Log
	trace   port.Trace
}

func NewObservability(
	metrics port.Metrics,
	log port.Log,
	trace port.Trace,
) *Observability {
	return &Observability{
		metrics: metrics,
		log:     log,
		trace:   trace,
	}
}

func (o *Observability) Metrics() port.Metrics {
	return o.metrics
}

func (o *Observability) Log() port.Log {
	return o.log
}

func (o *Observability) Trace() port.Trace {
	return o.trace
}
