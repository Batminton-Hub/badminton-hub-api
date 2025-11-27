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
	authBody := domain.AuthBody{}
	response := domain.RespAuth{}
	config := util.LoadConfig()

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, errInfo := memberUtil.ValidateBearerToken(authInfo.BearerToken)
	if errInfo.Err != nil {
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
		response.Resp = domain.ErrAuthenticateFail
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

type LoginSystem struct {
	ctx            context.Context
	memberRepo     port.MemberRepo
	middlewareUtil port.MiddlewareUtil
	thirdPartyUtil port.ThirdPartyUtil
	observability  port.Observability
}

func NewLoginSystem(
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
	thirdPartyUtil port.ThirdPartyUtil,
	observability port.Observability,
) *LoginSystem {
	return &LoginSystem{
		ctx:            ctx,
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
		thirdPartyUtil: thirdPartyUtil,
		observability:  observability,
	}
}

func (l *LoginSystem) Login(loginInfo domain.LoginInfo) (int, domain.RespLogin) {
	response := domain.RespLogin{}

	_, log, trace := Observe(l.observability)
	startTrace := trace.SetScope(loginInfo.ScopeName)
	traceCtx, err := trace.NewContext(loginInfo.TraceID, loginInfo.SpanID)
	if err != nil {
		response.Resp = domain.ErrInternalServer
		return response.Resp.HttpStatus, response
	}
	span := startTrace.CreateSpan(traceCtx, "login-nomal")
	defer span.End()

	config := util.LoadConfig()
	loginForm := loginInfo.LoginForm

	spanFindEmail := span.AddSpan("find-email-member")
	memberBody, errInfo := l.memberRepo.FindEmailMember(l.ctx, loginForm.Email)
	if errInfo.Err != nil {
		response.Resp = errInfo.Resp
		spanFindEmail.SetStatus(trace.Code().Error(), errInfo.Err.Error())
		spanFindEmail.End()

		file, line := GetLine()
		log.Error(l.ctx, domain.LogError{
			TraceID: loginInfo.TraceID,
			Path:    loginInfo.Path,
			Job:     loginInfo.Job,
			File:    file,
			Line:    line,
			Error:   errInfo.Err.Error(),
		})
		return http.StatusBadRequest, response
	}
	spanFindEmail.End()

	// Check password
	spanPassword := span.AddSpan("check-password")
	if memberBody.Password != HashPassword(loginForm.Password, config.KeyHashPassword) {
		response.Resp = domain.ErrLoginHashPassword
		spanPassword.SetStatus(trace.Code().Error(), response.Resp.Msg)
		spanPassword.SetTag(trace.Tag().String("hash_member_password", memberBody.Password))
		spanPassword.SetTag(trace.Tag().String("hash_login_password", HashPassword(loginForm.Password, config.KeyHashPassword)))
		spanPassword.End()

		file, line := GetLine()
		log.Error(l.ctx, domain.LogError{
			TraceID: loginInfo.TraceID,
			Path:    loginInfo.Path,
			Job:     loginInfo.Job,
			File:    file,
			Line:    line,
			Error:   response.Resp.Err.Error(),
		})

		return http.StatusUnauthorized, response
	}
	spanPassword.End()

	// create token
	spanToken := span.AddSpan("create-token")
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, errInfo := l.middlewareUtil.GenBearerToken(hashAuth)
	if errInfo.Err != nil {
		response.Resp = domain.ErrGenerateToken
		spanToken.SetStatus(trace.Code().Error(), errInfo.Err.Error())
		spanToken.End()

		file, line := GetLine()
		log.Error(l.ctx, domain.LogError{
			TraceID: loginInfo.TraceID,
			Path:    loginInfo.Path,
			Job:     loginInfo.Job,
			File:    file,
			Line:    line,
			Error:   response.Resp.Err.Error(),
		})
		return http.StatusInternalServerError, response
	}
	spanToken.End()

	response.BearerToken = tokenObj.Token
	response.Resp = domain.LoginSuccess

	logLogin := domain.LogInfo{
		TraceID: loginInfo.TraceID,
		Path:    loginInfo.Path,
		Job:     loginInfo.Job,
		Message: response.Resp.Msg,
	}
	l.observability.Log().Info(l.ctx, logLogin)

	return domain.LoginSuccess.HttpStatus, response
}

func (l *LoginSystem) LoginThirdParty(loginInfo domain.LoginInfo) (int, domain.RespLogin) {
	response := domain.RespLogin{}
	info, resp := l.thirdPartyUtil.BindingRequest(loginInfo.Platform, loginInfo.PlatformData)
	if resp.Status == domain.ERROR {
		response.Resp = resp
		return resp.HttpStatus, response
	}

	loginForm := domain.LoginForm{
		Email: info.Email,
	}
	memberBody, errInfo := l.memberRepo.FindEmailMember(l.ctx, loginForm.Email)
	if errInfo.Err != nil {
		response.Resp = errInfo.Resp
		return errInfo.Resp.HttpStatus, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, errInfo := l.middlewareUtil.GenBearerToken(hashAuth)
	if errInfo.Err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	response.Resp = domain.LoginSuccess
	return http.StatusOK, response
}

type RegisterSystem struct {
	Ctx            context.Context
	MemberRepo     port.MemberRepo
	MiddlewareUtil port.MiddlewareUtil
	ThirdPartyUtil port.ThirdPartyUtil
}

func NewRegisterSystem(
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
	thirdPartyUtil port.ThirdPartyUtil,
) *RegisterSystem {
	return &RegisterSystem{
		Ctx:            ctx,
		MemberRepo:     memberRepo,
		MiddlewareUtil: middlewareUtil,
		ThirdPartyUtil: thirdPartyUtil,
	}
}

func (r *RegisterSystem) Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister) {
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

	if errResult := r.MemberRepo.SaveMember(r.Ctx, memberBody); errResult.Err != nil {
		response.Resp = errResult.Resp
		return errResult.Resp.HttpStatus, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, errInfo := r.MiddlewareUtil.GenBearerToken(hashAuth)
	if errInfo.Err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.Resp = domain.RegisterSuccess
	response.BearerToken = tokenObj.Token
	return domain.RegisterSuccess.HttpStatus, response // 201 Created for successful registration
}

func (r *RegisterSystem) RegisterThirdParty(registerInfo domain.RegisterInfo) (int, domain.RespRegister) {
	response := domain.RespRegister{}
	info, resp := r.ThirdPartyUtil.BindingRequest(registerInfo.Platform, registerInfo.PlatformData)
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

	if errResult := r.MemberRepo.SaveMember(r.Ctx, memberBody); errResult.Err != nil {
		response.Resp = errResult.Resp
		return errResult.Resp.HttpStatus, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, errInfo := r.MiddlewareUtil.GenBearerToken(hashAuth)
	if errInfo.Err != nil {
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	response.Resp = domain.RedirectSuccess
	return http.StatusCreated, response
}

type MiddlewareUtil struct {
	encryption port.EncryptionUtil
}

func NewMiddlewareUtil(
	encryption port.EncryptionUtil,
) *MiddlewareUtil {
	return &MiddlewareUtil{
		encryption: encryption,
	}
}

// middleware util
func (m *MiddlewareUtil) ValidateBearerToken(tokenObj domain.BearerToken) (domain.AuthBody, domain.ErrInfo) {
	errInfo := domain.ErrInfo{}
	config := util.LoadConfig()

	token := tokenObj.Token[len("Bearer "):] // Remove "Bearer " prefix

	authBody := domain.AuthBody{}
	err := m.encryption.Decrypte(token, config.KeyBearerToken, &authBody)
	if err != nil {
		errInfo.Err = err
		errInfo.Resp = domain.ErrValidateToken
		return authBody, errInfo
	}

	if authBody.Exp < time.Now().Unix() {
		errInfo.Err = fmt.Errorf("token has expired")
		errInfo.Resp = domain.ErrValidateToken
		return authBody, errInfo
	}

	return authBody, errInfo
}

func (m *MiddlewareUtil) GenBearerToken(hashBody domain.HashAuth) (domain.BearerToken, domain.ErrInfo) {
	response := domain.BearerToken{}
	errInfo := domain.ErrInfo{}
	var token string
	config := util.LoadConfig()

	lt := config.BearerTokenExp
	createAt := time.Now().UTC()
	exp := time.Now().Add(lt).Unix()

	byteHash, err := util.EncryptGOB(hashBody)
	if err != nil {
		errInfo.Err = fmt.Errorf("failed to encrypt hash body: %w", err)
		errInfo.Resp = domain.ErrGenerateToken
		return response, errInfo
	}
	rawHash := string(byteHash)
	authBody := domain.AuthBody{
		CreateAt: createAt,
		Exp:      exp,
		Data: domain.AuthMember{
			UserID:    hashBody.UserID,
			CreatedAt: hashBody.CreateAt,
			HashAuth:  HashAuth(rawHash, config.KeyHashAuth),
		},
	}

	encryptedMember, err := m.encryption.Encrypte(authBody, config.KeyBearerToken, lt)
	if err != nil {
		errInfo.Err = fmt.Errorf("failed to encrypt member: %w", err)
		errInfo.Resp = domain.ErrGenerateToken
		return response, errInfo
	}

	token = encryptedMember
	response.Token = token
	return response, errInfo
}
