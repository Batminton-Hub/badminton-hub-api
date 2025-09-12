package port

import "Badminton-Hub/internal/core/domain"

type ThirdPartyUtil interface {
	BingingRequest(platform string, platformData any) (domain.ThirdPartyDataForm, domain.Resp)
}
