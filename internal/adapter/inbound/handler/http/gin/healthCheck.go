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

	tracer := h.observability.Trace()
	scopeName := "health-check"
	initTracer := tracer.InitTracer(scopeName)
	span := initTracer.ParaentSpan(c, "health-check")
	defer span.End()

	span.SetTag(tracer.Tag().String("test-string-key-1", "test-value-1"))
	span.SetTag(tracer.Tag().String("test-string-key-2", "test-value-2"))

	eventTag := []domain.TracerTag{
		tracer.Tag().String("test-string-key-3", "test-value-3"),
		tracer.Tag().String("test-string-key-4", "test-value-4"),
		tracer.Tag().Int64("test-int64-key-1", 123),
		tracer.Tag().Int64("test-int64-key-2", 456),
	}
	span.AddEvent("test-event-1", eventTag...)
	span.SetStatus(tracer.Code().OK(), "test-status-1")
	span.SetName("test-name-1")

	// spanID := span.GetSpan().ID()

	response := domain.HealthCheckSuccess
	Resp(c, response.HttpStatus, response.Code, response.Msg, nil)
}
