package port

import (
	"Badminton-Hub/internal/core/domain"
	"context"
)

type GangService interface {
	CreateGang(gang domain.Gang) (int, domain.RespCreateGang)
}

type GangRepo interface {
	SaveGang(ctx context.Context, gang domain.Gang) domain.ErrInfo
}
