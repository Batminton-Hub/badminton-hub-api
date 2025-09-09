package core_util

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"context"
	"net/http"
	"time"
)

type PlatformUtil interface {
	Login(domain.LoginInfo) (int, domain.RespLogin)
	Register(domain.RegisterInfo) (int, domain.RespRegister)
	Authenticate(domain.AuthInfo) (int, domain.RespAuth)
}

type NomalPlatform struct {
	CTX            context.Context
	MemberRepo     port.MemberRepo
	MiddlewareUtil port.MiddlewareUtil
}

type GooglePlatform struct {
	CTX            context.Context
	MemberRepo     port.MemberRepo
	MiddlewareUtil port.MiddlewareUtil
	Callback       port.CallbackService
}

func NewNomalPlatform(
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
) *NomalPlatform {
	return &NomalPlatform{
		CTX:            ctx,
		MemberRepo:     memberRepo,
		MiddlewareUtil: middlewareUtil,
	}
}

func NewGooglePlatform(
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
	callback port.CallbackService,
) *GooglePlatform {
	return &GooglePlatform{
		CTX:            ctx,
		MemberRepo:     memberRepo,
		MiddlewareUtil: middlewareUtil,
		Callback:       callback,
	}
}

func SwitchPlatform(
	platform string,
	ctx context.Context,
	memberRepo port.MemberRepo,
	middlewareUtil port.MiddlewareUtil,
) (PlatformUtil, error) {
	switch platform {
	case domain.NORMAL:
		return &NomalPlatform{ctx, memberRepo, middlewareUtil}, nil
	case domain.GOOGLE:
		return &GooglePlatform{ctx, memberRepo, middlewareUtil, nil}, nil
	}
	return nil, domain.ErrPlatformNotSupport.Err
}

func (n *NomalPlatform) Login(loginInfo domain.LoginInfo) (int, domain.RespLogin) {
	response := domain.RespLogin{}
	ctx := n.CTX

	config := util.LoadConfig()
	loginForm := loginInfo.LoginForm

	memberBody, err := n.MemberRepo.FindEmailMember(ctx, loginForm.Email)
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

	tokenObj, err := n.MiddlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	response.Resp = domain.LoginSuccess
	return domain.LoginSuccess.HttpStatus, response
}

func (g *GooglePlatform) Login(loginInfo domain.LoginInfo) (int, domain.RespLogin) {
	ctx := g.CTX
	response := domain.RespLogin{}
	info, ok := loginInfo.PlatformData.(domain.ResponseGoogleLoginCallback)
	if !ok {
		response.Resp = domain.ErrInvalidOAuthDecode
		return http.StatusBadRequest, response
	}

	loginForm := domain.LoginForm{
		Email: info.UserInfo.Email,
	}
	memberBody, err := g.MemberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Resp = domain.ErrMemberEmailNotFound
		return http.StatusBadRequest, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, err := g.MiddlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	return http.StatusOK, response
}

func (n *NomalPlatform) Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister) {
	ctx := n.CTX

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

	if err := n.MemberRepo.SaveMember(ctx, memberBody); err != nil {
		switch err {
		case domain.ErrMemberRegisterFailDuplicateEmail.Err:
			response.Resp = domain.ErrMemberRegisterFailDuplicateEmail
		case domain.ErrMemberRegisterFailDuplicateHash.Err:
			response.Resp = domain.ErrMemberRegisterFailDuplicateHash
		default:
			response.Resp = domain.ErrCreateMemberFail
		}
		return http.StatusInternalServerError, response // 500 Internal Server Error for other DB errors
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, err := n.MiddlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		response.Resp = domain.ErrGenerateToken
		return http.StatusInternalServerError, response
	}

	response.Resp = domain.RegisterSuccess
	response.BearerToken = tokenObj.Token
	return domain.RegisterSuccess.HttpStatus, response // 201 Created for successful registration
}

func (g *GooglePlatform) Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister) {
	ctx := g.CTX

	response := domain.RespRegister{}
	info, ok := registerInfo.PlatformData.(domain.ResponseGoogleRegisterCallback)
	if !ok {
		response.Resp = domain.ErrInvalidOAuthDecode
		return response.Resp.HttpStatus, response
	}

	config := util.LoadConfig()

	createAt := time.Now()
	updateAt := time.Now()
	password := util.Sha256(info.UserInfo.ID + domain.GOOGLE)
	memberBody := domain.Member{
		UserID:       GenUserID(info.UserInfo.Email, password),
		Email:        info.UserInfo.Email,
		Hash:         GenerateHash(config.KeyHashMember),
		DisplayName:  info.UserInfo.Name,
		ProfileImage: info.UserInfo.Picture,
		Status:       domain.PENDING,
		TypeMember:   domain.MEMBER,
		GoogleID:     info.UserInfo.ID,
		CreatedAt:    createAt,
		UpdatedAt:    updateAt,
	}
	if err := g.MemberRepo.SaveMember(ctx, memberBody); err != nil {
		response.Resp = domain.ErrCreateMemberFail
		return response.Resp.HttpStatus, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	tokenObj, err := g.MiddlewareUtil.GenBearerToken(hashAuth)
	if err != nil {
		return http.StatusInternalServerError, response
	}

	response.BearerToken = tokenObj.Token
	return http.StatusCreated, response
}

func (n *NomalPlatform) Authenticate(authInfo domain.AuthInfo) (int, domain.RespAuth) {
	response := domain.RespAuth{}
	config := util.LoadConfig()

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := n.MiddlewareUtil.ValidateBearerToken(authInfo.BearerToken)
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

func (g *GooglePlatform) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	return g.Callback.Authenticate(info)
}
