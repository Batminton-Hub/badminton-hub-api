package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type HealthCheckController interface {
	HealthCheck(c *gin.Context)
}

type HealthCheck struct {
	observability port.Observability
}

func (h *HealthCheck) HealthCheck(c *gin.Context) {
	counter := domain.MetricsCounter{
		Name: "health-check-request",
		Help: "Number of health check requests",
	}
	countHealthCheck := h.observability.Metrics().Counter(counter)
	countHealthCheck.Inc()

	logInfo := domain.LogInfo{
		Path:    getPath(c),
		Job:     "health-check",
		Message: "Health check success",
	}
	h.observability.Log().Info(c, logInfo)

	response := domain.HealthCheckSuccess
	Resp(c, response.HttpStatus, response.Code, response.Msg, nil)
}
