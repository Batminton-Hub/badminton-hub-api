package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type ObservabilityController interface {
	Metrics(c *gin.Context)
}

type Observability struct {
	observability port.Observability
}

func (o *Observability) Metrics(c *gin.Context) {
	info := domain.MetricsHttp{
		Writer:  c.Writer,
		Request: c.Request,
	}
	o.observability.Metrics().GetMetrics(info)
}
