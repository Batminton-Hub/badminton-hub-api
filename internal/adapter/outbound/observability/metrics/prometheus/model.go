package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Prometheus struct {
	CounterInfo   map[string]prometheus.Counter
	GaugeInfo     map[string]prometheus.Gauge
	HistogramInfo map[string]prometheus.Histogram
}

type CounterServiceImpl struct {
	Counter prometheus.Counter
}

type GaugeServiceImpl struct {
	Gauge prometheus.Gauge
}

type HistogramServiceImpl struct {
	Histogram prometheus.Histogram
}
