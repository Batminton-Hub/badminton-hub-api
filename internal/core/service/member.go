package service

// type MemberService struct {
// 	memberRepo     port.MemberRepo
// 	middlewareUtil port.MiddlewareService
// }

// type MemberServiceGoogle struct {
// 	memberRepo     port.MemberRepo
// 	middlewareUtil port.MiddlewareService
// }

// func NewMemberService(memberRepo port.MemberRepo, middlewareUtil port.MiddlewareService) *MemberService {
// 	memberService := &MemberService{
// 		memberRepo:     memberRepo,
// 		middlewareUtil: middlewareUtil,
// 	}
// 	return memberService
// }

// func (m *MemberService) GetProfile(userID string) (int, domain.ResponseGetProfile) {
// 	ctx, cancel := util.InitConText(2 * time.Second)
// 	defer cancel()

// 	response := domain.ResponseGetProfile{}
// 	member, err := m.memberRepo.GetMemberByUserID(ctx, userID)
// 	if err != nil {
// 		response.Code = domain.ErrGetMember.Code
// 		response.Message = domain.ErrGetMember.Msg
// 		return http.StatusBadRequest, response
// 	}

// 	response.Code = domain.Success.Code
// 	response.Message = domain.Success.Msg
// 	response.Member = member
// 	return http.StatusOK, response
// }

// func (m *MemberService) UpdateProfile(userID string, request domain.RequestUpdateProfile) (int, domain.ResponseUpdateProfile) {
// 	ctx, cancel := util.InitConText(2 * time.Second)
// 	defer cancel()

// 	response := domain.ResponseUpdateProfile{}

// 	if request.DisplayName == "" &&
// 		request.ProfileImage == "" &&
// 		request.DateOfBirth == "" &&
// 		request.Region == "" &&
// 		request.Gender == "" &&
// 		request.Phone == "" &&
// 		len(request.Tag) == 0 {
// 		response.Code = domain.ErrInvalidInput.Code
// 		response.Message = domain.ErrInvalidInput.Msg
// 		return http.StatusBadRequest, response
// 	}

// 	if err := m.memberRepo.UpdateMember(ctx, userID, request); err != nil {
// 		response.Code = domain.ErrUpdateMemberFail.Code
// 		response.Message = domain.ErrUpdateMemberFail.Msg
// 		return http.StatusInternalServerError, response
// 	}

// 	response.Code = domain.UpdateMemberSuccess.Code
// 	response.Message = domain.UpdateMemberSuccess.Msg
// 	return http.StatusOK, response
// }
