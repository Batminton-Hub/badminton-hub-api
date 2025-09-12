package google

import "Badminton-Hub/internal/core/domain"

func BingingRequest(platformData any) (domain.ThirdPartyDataForm, domain.Resp) {
	response := domain.ThirdPartyDataForm{}
	info, ok := platformData.(GoogleMemberInfo)
	if !ok {
		return response, domain.ErrInvalidDecode3rdPartyForm
	}
	response.Platform = domain.GOOGLE
	response.Email = info.UserInfo.Email
	response.DisplayName = info.UserInfo.Name
	response.PlatformID = info.UserInfo.ID
	return response, domain.Success
}
