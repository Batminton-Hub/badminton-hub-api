package third_party

import (
	"Badminton-Hub/internal/adapter/outbound/3rdParty/google"
	"Badminton-Hub/internal/core/domain"
)

type ThirdPartyUtilImpl struct {
}

func NewThirdPartyUtil() *ThirdPartyUtilImpl {
	return &ThirdPartyUtilImpl{}
}

func (t *ThirdPartyUtilImpl) BindingRequest(platform string, platformData any) (domain.ThirdPartyDataForm, domain.Resp) {
	dataForm := domain.ThirdPartyDataForm{}
	switch platform {
	case domain.GOOGLE:
		return google.BindingRequest(platformData)
	default:
		return dataForm, domain.ErrPlatformNotSupport
	}
}
