package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"fmt"

	"github.com/gin-gonic/gin"
)

type HealthCheckController interface {
	HealthCheck(c *gin.Context)
}

type HealthCheck struct {
	observability port.Observability
}

func (h *HealthCheck) HealthCheck(c *gin.Context) {
	// counter := domain.MetricsCounter{
	// 	Name: "health-check-request",
	// 	Help: "Number of health check requests",
	// }
	// countHealthCheck := h.observability.Metrics().Counter(counter)
	// countHealthCheck.Inc()

	// logInfo := domain.LogInfo{
	// 	Path:    getPath(c),
	// 	Job:     "health-check",
	// 	Message: "Health check success",
	// }
	// h.observability.Log().Info(c, logInfo)

	// scopeName := "health-check"
	// trace := h.observability.Trace().SetScope(scopeName)
	// span := trace.CreateSpan(c, "start-trace")
	// defer span.End()
	// fmt.Println("traceID :", span.GetTraceID())

	// span1 := span.AddSpan("health-check-span-1")
	// span1.End()

	// span2 := span1.AddSpan("health-check-span-2")
	// span2.End()

	// TestHealthCheck(h.observability, span.GetTraceID(), span2.GetSpanID())

	response := domain.HealthCheckSuccess
	Resp(c, response.HttpStatus, response.Code, response.Msg, nil)
}

func TestHealthCheck(observability port.Observability, traceID string, spanID string) {
	newCtx, err := observability.Trace().NewContext(traceID, spanID)
	if err != nil {
		fmt.Println("error :", err)
	}
	scopeName := "health-check"
	trace := observability.Trace().SetScope(scopeName)
	newSpan := trace.ConnectSpan(newCtx)

	span1 := newSpan.AddSpan("health-check-span-3")
	span1.End()

	span2 := newSpan.AddSpan("health-check-span-4")
	span2.End()

	fmt.Println("traceID in function:", span1.GetTraceID())
}
