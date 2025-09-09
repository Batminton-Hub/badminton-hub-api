package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/internal/core_util"
	"Badminton-Hub/util"
	"fmt"
	"time"
)

type AuthenticationService struct {
	memberRepo     port.MemberRepo
	encryption     port.EncryptionUtil
	middlewareUtil port.MiddlewareUtil
}

type MiddlewareService struct {
	memberRepo port.MemberRepo
	encryption port.EncryptionUtil
	callback   port.CallbackService
}

type AuthenticationSystem struct {
	port.AuthenticationService
	port.MiddlewareService
}

func NewAuthenticationSystem(
	authenticationService port.AuthenticationService,
	middlewareService port.MiddlewareService,
) *AuthenticationSystem {
	return &AuthenticationSystem{
		authenticationService,
		middlewareService,
	}
}

func NewAuthenticationService(
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
) *AuthenticationService {
	return &AuthenticationService{
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
	}
}

func NewMiddlewareService(
	memberRepo port.MemberRepo,
	encryption port.EncryptionUtil,
	callback *CallbackService,
) *MiddlewareService {
	return &MiddlewareService{
		memberRepo: memberRepo,
		encryption: encryption,
		callback:   callback,
	}
}

func (a *AuthenticationService) Login(loginInfo domain.LoginInfo) (int, domain.RespLogin) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespLogin{}
	platformUtil, err := core_util.SwitchPlatform(loginInfo.Platform, ctx, a.memberRepo, a.middlewareUtil)
	if err != nil {
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}

	return platformUtil.Login(loginInfo)
}

func (a *AuthenticationService) Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespRegister{}
	platformUtil, err := core_util.SwitchPlatform(registerInfo.Platform, ctx, a.memberRepo, a.middlewareUtil)
	if err != nil {
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}

	return platformUtil.Register(registerInfo)
}

func (m *MiddlewareService) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	var response domain.RespAuth
	var platformUtil core_util.PlatformUtil
	platform := info.Platform
	switch platform {
	case domain.NORMAL:
		platformUtil = core_util.NewNomalPlatform(ctx, m.memberRepo, m)
	case domain.GOOGLE:
		platformUtil = core_util.NewGooglePlatform(ctx, m.memberRepo, m, m.callback)
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return domain.ErrPlatformNotSupport.HttpStatus, response
	}

	return platformUtil.Authenticate(info)
}

func (m *MiddlewareService) ValidateBearerToken(tokenObj domain.BearerToken) (domain.AuthBody, error) {
	config := util.LoadConfig()

	token := tokenObj.Token[len("Bearer "):] // Remove "Bearer " prefix

	authBody := domain.AuthBody{}
	err := m.encryption.Decrypte(token, config.KeyBearerToken, &authBody)
	if err != nil {
		return authBody, err
	}

	if authBody.Exp < time.Now().Unix() {
		fmt.Println("authBody.Exp < time.Now().Unix() : ", authBody.Exp, " : ", time.Now().Unix())
		return authBody, fmt.Errorf("token has expired")
	}

	return authBody, nil
}

func (m *MiddlewareService) GenBearerToken(hashBody domain.HashAuth) (domain.BearerToken, error) {
	response := domain.BearerToken{}
	var token string
	config := util.LoadConfig()

	lt := config.BearerTokenExp
	createAt := time.Now().UTC()
	exp := time.Now().Add(lt).Unix()

	byteHash, err := util.EncryptGOB(hashBody)
	if err != nil {
		return response, fmt.Errorf("failed to encrypt hash body: %w", err)
	}
	rawHash := string(byteHash)
	authBody := domain.AuthBody{
		CreateAt: createAt,
		Exp:      exp,
		Data: domain.AuthMember{
			UserID:    hashBody.UserID,
			CreatedAt: hashBody.CreateAt,
			HashAuth:  core_util.HashAuth(rawHash, config.KeyHashAuth),
		},
	}

	encryptedMember, err := m.encryption.Encrypte(authBody, config.KeyBearerToken, lt)
	if err != nil {
		return response, fmt.Errorf("failed to encrypt member: %w", err)
	}

	token = encryptedMember
	response.Token = token
	return response, nil
}
