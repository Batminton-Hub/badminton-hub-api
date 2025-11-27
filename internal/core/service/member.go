package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"time"
)

type MemberService struct {
	memberRepo     port.MemberRepo
	middlewareUtil port.MiddlewareService
}

func NewMemberService(
	memberRepo port.MemberRepo,
) *MemberService {
	memberService := &MemberService{
		memberRepo: memberRepo,
	}
	return memberService
}

func (m *MemberService) GetProfile(userInfo domain.ReqGetProfile) (int, domain.RespGetProfile) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespGetProfile{}
	member, errInfo := m.memberRepo.GetMemberByUserID(ctx, userInfo.UserID)
	if errInfo.Err != nil {
		response.Resp = domain.ErrGetMember
		return response.Resp.HttpStatus, response
	}

	response.Member = member
	response.Resp = domain.Success
	return response.Resp.HttpStatus, response
}

func (m *MemberService) UpdateProfile(userInfo domain.ReqGetProfile, request domain.ReqUpdateProfile) (int, domain.RespUpdateProfile) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespUpdateProfile{}
	if request.DisplayName == "" &&
		request.ProfileImage == "" &&
		request.DateOfBirth == "" &&
		request.Region == "" &&
		request.Gender == "" &&
		request.Phone == "" &&
		len(request.Tag) == 0 {
		response.Resp = domain.ErrInvalidInput
		return response.Resp.HttpStatus, response
	}

	userID := userInfo.UserID
	if errInfo := m.memberRepo.UpdateMember(ctx, userID, request); errInfo.Err != nil {
		response.Resp = domain.ErrUpdateMemberFail
		return response.Resp.HttpStatus, response
	}

	response.Resp = domain.UpdateMemberSuccess
	return response.Resp.HttpStatus, response
}
