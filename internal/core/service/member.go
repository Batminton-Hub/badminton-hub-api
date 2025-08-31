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

	createAt := time.Now()
	updateAt := time.Now()
	memberBody := domain.Member{
		UserID:    util.GenUserID(registerForm.Email, registerForm.Password),
		Email:     registerForm.Email,
		Password:  util.HashPassword(registerForm.Password, config.KeyHashPassword),
		Gender:    registerForm.Gender,
		Hash:      util.GenerateHash(config.KeyHashMember),
		CreatedAt: createAt,
		UpdatedAt: updateAt,
	}

	if err := m.memberRepo.RegisterMember(ctx, memberBody); err != nil {
		response.ErrorCode = domain.ErrCreateMemberFail.Code
		response.Error = domain.ErrCreateMemberFail.Err
		response.Message = domain.ErrCreateMemberFail.Msg
		return 400, response
	}

	// create token
	hashAuth := domain.HashAuth{
		Username: memberBody.Username,
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := util.GenBearerToken(hashAuth, m.middlewareUtil.Encryptetion())
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
	response := domain.ResponseLogin{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	config, err := util.LoadConfig()
	if err != nil {
		response.ErrorCode = domain.ErrLoadConfig.Code
		response.Error = domain.ErrLoadConfig.Err
		return 500, response
	}

	memberBody, err := m.memberRepo.FindEmailMember(ctx, loginForm)
	if err != nil {
		response.ErrorCode = domain.ErrMemberEmailNotFound.Code
		response.Error = domain.ErrMemberEmailNotFound.Err
		return 400, response
	}

	// Check password
	if memberBody.Password != util.HashPassword(loginForm.Password, config.KeyHashPassword) {
		response.ErrorCode = domain.ErrLoginHashPassword.Code
		response.Error = domain.ErrLoginHashPassword.Err
		return 401, response
	}

	// create token
	hashAuth := domain.HashAuth{
		Username: memberBody.Username,
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := util.GenBearerToken(hashAuth, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.ErrorCode = domain.ErrGenerateToken.Code
		response.Error = domain.ErrGenerateToken.Err
		return 500, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return 200, response
}
