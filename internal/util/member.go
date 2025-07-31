package core_util

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"context"
	"time"
)

type MemberUtil struct {
	memberRepo port.MemberRepo
}

func NewMemberUtil(memberRepo port.MemberRepo) *MemberUtil {
	return &MemberUtil{
		memberRepo: memberRepo,
	}
}
func (m *MemberUtil) RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response := domain.ResponseRegisterMember{}
	memberBody := domain.Member{
		Email:     registerForm.Email,
		Password:  util.HashPassword(registerForm.Password, "test"),
		Gender:    registerForm.Gender,
		Hash:      util.GenerateHash(registerForm.Email),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := m.memberRepo.RegisterMember(ctx, memberBody); err != nil {
		response.ErrorCode = 1001
		response.Error = "Failed to register member"
		return 400, response
	}

	// create token
	var token string
	token, err := util.GenBearerToken(memberBody)
	if err != nil {
		response.ErrorCode = 1002
		response.Error = "Failed to generate token"
		return 500, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return 200, response
}

func (m *MemberUtil) Login(loginForm domain.LoginForm) (int, domain.ResponseLogin) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response := domain.ResponseLogin{}
	memberBody, err := m.memberRepo.LoginByEmail(ctx, loginForm)
	if err != nil {
		response.ErrorCode = 1003
		response.Error = "Failed to login member"
		return 400, response
	}

	// Check password
	if memberBody.Password != util.HashPassword(loginForm.Password, "test") {
		response.ErrorCode = 1004
		response.Error = "Invalid email or password"
		return 401, response
	}

	// create token
	var token string
	token, err = util.GenBearerToken(memberBody)
	if err != nil {
		response.ErrorCode = 1005
		response.Error = "Failed to generate token"
		return 500, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return 200, response
}
