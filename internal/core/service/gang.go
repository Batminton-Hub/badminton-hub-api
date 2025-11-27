package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"time"
)

type GangService struct {
	observability port.Observability
	gangRepo      port.GangRepo
}

func NewGangService(
	gangRepo port.GangRepo,
	observability port.Observability,
) *GangService {
	return &GangService{
		gangRepo:      gangRepo,
		observability: observability,
	}
}

func (g *GangService) CreateGang(gang domain.Gang) (int, domain.RespCreateGang) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespCreateGang{}

	timeNow := time.Now()
	gang.GangID = util.GenerateUUID()
	gang.CreatedAt = timeNow
	gang.UpdatedAt = timeNow

	errInfo := g.gangRepo.SaveGang(ctx, gang)
	if errInfo.Err != nil {
		response.Resp = domain.ErrCreateGangFail
		return response.Resp.HttpStatus, response
	}

	response.GangID = gang.GangID
	response.Resp = domain.Success
	return response.Resp.HttpStatus, response
}
