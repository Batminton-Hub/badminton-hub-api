package prometheus

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheus() *Prometheus {
	return &Prometheus{
		CounterInfo:   map[string]prometheus.Counter{},
		GaugeInfo:     map[string]prometheus.Gauge{},
		HistogramInfo: map[string]prometheus.Histogram{},
	}
}

func (p *Prometheus) GetMetrics(info domain.MetricsHttp) {
	promhttp.Handler().ServeHTTP(info.Writer, info.Request)
}

func (p *Prometheus) Counter(info domain.MetricsCounter) port.MetricsCounterUtil {
	if _, ok := p.CounterInfo[info.Name]; !ok {
		p.CounterInfo[info.Name] = promauto.NewCounter(prometheus.CounterOpts{
			Name: info.Name,
			Help: info.Help,
		})
	}
	return &CounterServiceImpl{
		Counter: p.CounterInfo[info.Name],
	}
}

func (p *Prometheus) Gauge(info domain.MetricsGauge) port.MetricsGaugeUtil {
	if _, ok := p.GaugeInfo[info.Name]; !ok {
		p.GaugeInfo[info.Name] = promauto.NewGauge(prometheus.GaugeOpts{
			Name: info.Name,
			Help: info.Help,
		})
	}
	return &GaugeServiceImpl{
		Gauge: p.GaugeInfo[info.Name],
	}
}

func (c *CounterServiceImpl) Inc() {
	c.Counter.Inc()
}
func (c *CounterServiceImpl) Add(value float64) error {
	if value < 0 {
		return errors.New("value must be greater than 0")
	}
	c.Counter.Add(value)
	return nil
}

func (g *GaugeServiceImpl) Set(value float64) {
	g.Gauge.Set(value)
}
func (g *GaugeServiceImpl) Inc() {
	g.Gauge.Inc()
}
func (g *GaugeServiceImpl) Dec() {
	g.Gauge.Dec()
}
func (g *GaugeServiceImpl) Add(value float64) {
	g.Gauge.Add(value)
}
func (g *GaugeServiceImpl) Sub(value float64) {
	g.Gauge.Sub(value)
}

func (h *HistogramServiceImpl) Observe(value float64) {
	h.Histogram.Observe(value)
}
