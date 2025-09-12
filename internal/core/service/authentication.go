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
	thirdPartyUtil port.ThirdPartyUtil
	memberRepo     port.MemberRepo
	encryption     port.EncryptionUtil
	middlewareUtil port.MiddlewareUtil
}

type AuthenticateService struct {
	thirdParty     port.AuthenticateUtil
	middlewareUtil port.MiddlewareUtil
	memberRepo     port.MemberRepo
}

type MiddlewareUtil struct {
	encryption port.EncryptionUtil
}

type AuthenticationSystem struct {
	port.AuthenticationService
	port.MiddlewareService
}

type MiddlewareSystem struct {
	port.AuthenticateUtil
	port.MiddlewareUtil
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
	thirdPartyUtil port.ThirdPartyUtil,
) *AuthenticationService {
	return &AuthenticationService{
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
		thirdPartyUtil: thirdPartyUtil,
	}
}

func NewAuthenticateService(
	thirdParty port.AuthenticateUtil,
	middlewareUtil port.MiddlewareUtil,
	memberRepo port.MemberRepo,
) *AuthenticateService {
	return &AuthenticateService{
		thirdParty:     thirdParty,
		middlewareUtil: middlewareUtil,
		memberRepo:     memberRepo,
	}
}

func NewMiddlewareUtil(
	encryption port.EncryptionUtil,
) *MiddlewareUtil {
	return &MiddlewareUtil{
		encryption: encryption,
	}
}

func NewMiddlewareSystem(
	middlewareService port.AuthenticateUtil,
	middlewareUtil port.MiddlewareUtil,
) *MiddlewareSystem {
	return &MiddlewareSystem{
		middlewareService,
		middlewareUtil,
	}
}

func (a *AuthenticationService) Login(loginInfo domain.LoginInfo) (int, domain.RespLogin) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespLogin{}
	switch loginInfo.TypeSystem {
	case domain.SYSTEM:
		return core_util.Login(loginInfo, ctx, a.memberRepo, a.middlewareUtil)
	case domain.THIRD_PARTY:
		return core_util.LoginThirdParty(loginInfo, ctx, a.memberRepo, a.middlewareUtil, a.thirdPartyUtil)
	default:
		response.Resp = domain.ErrSystemNotSupport
		return response.Resp.HttpStatus, response
	}
}

func (a *AuthenticationService) Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.RespRegister{}
	switch registerInfo.TypeSystem {
	case domain.SYSTEM:
		return core_util.Register(registerInfo, ctx, a.memberRepo, a.middlewareUtil)
	case domain.THIRD_PARTY:
		return core_util.RegisterThirdParty(registerInfo, ctx, a.memberRepo, a.middlewareUtil, a.thirdPartyUtil)
	default:
		response.Resp = domain.ErrSystemNotSupport
		return response.Resp.HttpStatus, response
	}
}

func (m *AuthenticateService) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	var response domain.RespAuth
	switch info.TypeSystem {
	case domain.SYSTEM:
		return core_util.Authenticate(info, m.middlewareUtil)
	case domain.THIRD_PARTY:
		return m.thirdParty.Authenticate(info)
	default:
		response.Resp = domain.ErrSystemNotSupport
		return response.Resp.HttpStatus, response
	}
}

func (m *MiddlewareUtil) ValidateBearerToken(tokenObj domain.BearerToken) (domain.AuthBody, error) {
	config := util.LoadConfig()

	token := tokenObj.Token[len("Bearer "):] // Remove "Bearer " prefix

	authBody := domain.AuthBody{}
	err := m.encryption.Decrypte(token, config.KeyBearerToken, &authBody)
	if err != nil {
		return authBody, err
	}

	if authBody.Exp < time.Now().Unix() {
		return authBody, fmt.Errorf("token has expired")
	}

	return authBody, nil
}

func (m *MiddlewareUtil) GenBearerToken(hashBody domain.HashAuth) (domain.BearerToken, error) {
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
