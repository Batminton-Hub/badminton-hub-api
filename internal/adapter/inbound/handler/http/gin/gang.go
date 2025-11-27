package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type GangController interface {
	CreateGang(c *gin.Context)
}

type Gang struct {
	observability port.Observability
	gangService   port.GangService
}

func (g *Gang) CreateGang(c *gin.Context) {
	observe := g.observability
	metrics := observe.Metrics()
	log := observe.Log()
	trace := observe.Trace()

	counter := metrics.Counter(domain.MetricsCounter{
		Name: "create-gang-request",
		Help: "Number of create gang requests",
	})

	scrope := trace.SetScope("create-gang")
	traceID := scrope.CreateSpan(c, "create-gang").GetTraceID()

	logInfo := domain.LogInfo{
		TraceID: traceID,
		Path:    getPath(c),
		Job:     "create-gang",
		Message: "Create gang request received",
	}

	logError := domain.LogError{
		TraceID: traceID,
		Path:    getPath(c),
		Job:     "create-gang",
	}

	gang := domain.Gang{}
	if err := c.ShouldBindJSON(&gang); err != nil {
		file, line := observe.GetLine()
		logError.File = file
		logError.Line = line
		log.Error(c, logError)
		Resp(c, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg, nil)
		return
	}

	userID := getUserID(c)
	gang.OwnerID = userID
	httpStatus, response := g.gangService.CreateGang(gang)
	if response.Resp.Status == domain.ERROR {
		file, line := observe.GetLine()
		logError.File = file
		logError.Line = line
		log.Error(c, logError)
		Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, nil)
		return
	}

	counter.Inc()
	log.Info(c, logInfo)

	Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, response.GangID)
}
