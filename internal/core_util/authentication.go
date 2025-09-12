package core_util

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func HashAuth(rawHash, key string) string {
	data := fmt.Sprint(rawHash + key)
	hashAuth := util.Sha256(data)
	return hashAuth
}

func RandomGoogleState() (string, error) {
	config := util.LoadConfig()

	if strings.Contains(domain.DEVERLOP, config.Mode) {
		return config.DefaultGoogleState, nil
	}

	state := util.RandomString(32, true, true, false)
	return string(state), nil
}

func Authenticate(authInfo domain.AuthInfo, memberUtil port.MiddlewareUtil) (int, domain.RespAuth) {
	response := domain.RespAuth{}
	config := util.LoadConfig()

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := memberUtil.ValidateBearerToken(authInfo.BearerToken)
	if err != nil {
		response.Resp = domain.ErrValidateToken
		return http.StatusUnauthorized, response
	}

	response.AuthBody = authBody

	// ตรวจสอบความถูกต้องของ token
	hashAuthBody := domain.HashAuth{
		CreateAt: authBody.Data.CreatedAt,
		UserID:   authBody.Data.UserID,
	}
	byteHash, err := util.EncryptGOB(hashAuthBody)
	if err != nil {
		return http.StatusUnauthorized, response
	}
	rawHash := string(byteHash)
	hashauth := HashAuth(rawHash, config.KeyHashAuth)
	if authBody.Data.HashAuth != hashauth {
		response.Resp = domain.ErrValidateHashAuth
		return http.StatusUnauthorized, response
	}

	// ตรวจสอบว่า token ยังไม่หมดอายุ
	if authBody.Exp < time.Now().Unix() {
		response.Resp = domain.ErrTokenExpired
		return http.StatusUnauthorized, response
	}

	response.Resp = domain.AuthSuccess
	return domain.AuthSuccess.HttpStatus, response
}

func Login(loginInfo domain.LoginInfo, ctx context.Context, memberRepo port.MemberRepo, middlewareUtil port.MiddlewareUtil) (int, domain.RespLogin) {
	response := domain.RespLogin{}

	config := util.LoadConfig()
	loginForm := loginInfo.LoginForm

	memberBody, err := memberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Resp = domain.ErrMemberEmailNotFound
		return http.StatusBadRequest, response
	}

	// Check password
	if memberBody.Password != HashPassword(loginForm.Password, config.KeyHashPassword) {
		response.Resp = domain.ErrLoginHashPassword
		return http.StatusUnauthorized, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}

	tokenObj, err := middlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	response.Resp = domain.LoginSuccess
	return domain.LoginSuccess.HttpStatus, response
}

func Register(registerInfo domain.RegisterInfo, ctx context.Context, memberRepo port.MemberRepo, middlewareUtil port.MiddlewareUtil) (int, domain.RespRegister) {
	registerForm := registerInfo.RegisterForm
	response := domain.RespRegister{}
	config := util.LoadConfig()

	createAt := time.Now()
	updateAt := time.Now()
	memberBody := domain.Member{
		UserID:     GenUserID(registerForm.Email, registerForm.Password),
		Email:      registerForm.Email,
		Password:   HashPassword(registerForm.Password, config.KeyHashPassword),
		Gender:     registerForm.Gender,
		Hash:       GenerateHash(config.KeyHashMember),
		Status:     domain.ACTIVE,
		TypeMember: domain.MEMBER,
		CreatedAt:  createAt,
		UpdatedAt:  updateAt,
	}

	if err := memberRepo.SaveMember(ctx, memberBody); err != nil {
		switch err {
		case domain.ErrMemberRegisterFailDuplicateEmail.Err:
			response.Resp = domain.ErrMemberRegisterFailDuplicateEmail
		case domain.ErrMemberRegisterFailDuplicateHash.Err:
			response.Resp = domain.ErrMemberRegisterFailDuplicateHash
		default:
			response.Resp = domain.ErrCreateMemberFail
		}
		return http.StatusInternalServerError, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, err := middlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.Resp = domain.RegisterSuccess
	response.BearerToken = tokenObj.Token
	return domain.RegisterSuccess.HttpStatus, response // 201 Created for successful registration
}

func LoginThirdParty(
	loginInfo domain.LoginInfo,
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
	thirdPartyUtil port.ThirdPartyUtil,
) (int, domain.RespLogin) {
	response := domain.RespLogin{}
	info, resp := thirdPartyUtil.BingingRequest(loginInfo.Platform, loginInfo.PlatformData)
	if resp.Status == domain.ERROR {
		response.Resp = resp
		return resp.HttpStatus, response
	}

	loginForm := domain.LoginForm{
		Email: info.Email,
	}
	memberBody, err := memberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Resp = domain.ErrMemberEmailNotFound
		return http.StatusBadRequest, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, err := middlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	response.Resp = domain.LoginSuccess
	return http.StatusOK, response
}

func RegisterThirdParty(
	registerInfo domain.RegisterInfo,
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
	thirdPartyUtil port.ThirdPartyUtil,
) (int, domain.RespRegister) {
	response := domain.RespRegister{}
	info, resp := thirdPartyUtil.BingingRequest(registerInfo.Platform, registerInfo.PlatformData)
	if resp.Status == domain.ERROR {
		response.Resp = resp
		return resp.HttpStatus, response
	}

	config := util.LoadConfig()

	createAt := time.Now()
	updateAt := time.Now()
	password := util.Sha256(info.PlatformID + domain.GOOGLE)
	memberBody := domain.Member{
		UserID:       GenUserID(info.Email, password),
		Email:        info.Email,
		Hash:         GenerateHash(config.KeyHashMember),
		DisplayName:  info.DisplayName,
		ProfileImage: info.Picture,
		Status:       domain.PENDING,
		TypeMember:   domain.MEMBER,
		GoogleID:     info.PlatformID,
		CreatedAt:    createAt,
		UpdatedAt:    updateAt,
	}
	if err := memberRepo.SaveMember(ctx, memberBody); err != nil {
		switch err {
		case domain.ErrMemberRegisterFailDuplicateEmail.Err:
			response.Resp = domain.ErrMemberRegisterFailDuplicateEmail
		case domain.ErrMemberRegisterFailDuplicateHash.Err:
			response.Resp = domain.ErrMemberRegisterFailDuplicateHash
		default:
			response.Resp = domain.ErrCreateMemberFail
		}
		return http.StatusInternalServerError, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, err := middlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	response.Resp = domain.RedirectSuccess
	return http.StatusCreated, response
}
