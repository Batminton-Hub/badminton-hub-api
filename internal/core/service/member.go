package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"context"
	"time"
)

type MemberUtil struct {
	memberRepo     port.MemberRepo
	middlewareUtil port.MiddlewareUtil
}

func NewMemberUtil(memberRepo port.MemberRepo, middlewareUtil port.MiddlewareUtil) *MemberUtil {
	return &MemberUtil{
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
	}
}
func (m *MemberUtil) RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response := domain.ResponseRegisterMember{}
	config, err := util.LoadConfig()
	if err != nil {
		response.ErrorCode = domain.ErrLoadConfig.Code
		response.Error = domain.ErrLoadConfig.Err
		response.Message = domain.ErrLoadConfig.Msg
		return 500, response
	}
	memberBody := domain.Member{
		Email:     registerForm.Email,
		Password:  util.HashPassword(registerForm.Password, config.KeyHashPassword),
		Gender:    registerForm.Gender,
		Hash:      util.GenerateHash(registerForm.Email),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := m.memberRepo.RegisterMember(ctx, memberBody); err != nil {
		response.ErrorCode = domain.ErrCreateMemberFail.Code
		response.Error = domain.ErrCreateMemberFail.Err
		response.Message = domain.ErrCreateMemberFail.Msg
		return 400, response
	}

	// create token
	token, err := util.GenBearerToken(memberBody, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.ErrorCode = domain.ErrGenerateToken.Code
		response.Error = domain.ErrGenerateToken.Err
		response.Message = domain.ErrGenerateToken.Msg
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
	token, err = util.GenBearerToken(memberBody, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.ErrorCode = 1005
		response.Error = "Failed to generate token"
		return 500, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return 200, response
}
