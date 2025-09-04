package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	core_util "Badminton-Hub/internal/util"
	"Badminton-Hub/util"
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	register = "REGISTER"
	login    = "LOGIN"
	google   = "GOOGLE"
)

type MemberUtil struct {
	memberRepo     port.MemberRepo
	middlewareUtil port.MiddlewareUtil
}

func NewMemberUtil(memberRepo port.MemberRepo, middlewareUtil port.MiddlewareUtil) *MemberUtil {
	memberUtil := &MemberUtil{
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
	}
	return memberUtil
}

func (m *MemberUtil) RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response := domain.ResponseRegisterMember{}
	config, err := util.LoadConfig()
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response // 500 Internal Server Error for config load fail
	}

	createAt := time.Now()
	updateAt := time.Now()
	memberBody := domain.Member{
		UserID:    core_util.GenUserID(registerForm.Email, registerForm.Password),
		Email:     registerForm.Email,
		Password:  core_util.HashPassword(registerForm.Password, config.KeyHashPassword),
		Gender:    registerForm.Gender,
		Hash:      core_util.GenerateHash(config.KeyHashMember),
		CreatedAt: createAt,
		UpdatedAt: updateAt,
	}

	if err := m.memberRepo.RegisterMember(ctx, memberBody); err != nil {
		if domain.ErrMemberRegisterFailDuplicateHash.Err == err {
			response.Code = domain.ErrMemberRegisterFailDuplicateHash.Code
			response.Message = domain.ErrMemberRegisterFailDuplicateHash.Msg
			return http.StatusConflict, response // 409 Conflict
		}
		response.Code = domain.ErrCreateMemberFail.Code
		response.Message = domain.ErrCreateMemberFail.Msg
		return http.StatusInternalServerError, response // 500 Internal Server Error for other DB errors
	}

	// create token
	hashAuth := domain.HashAuth{
		Username: memberBody.Username,
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return http.StatusCreated, response // 201 Created for successful registration
}

func (m *MemberUtil) Login(loginForm domain.LoginForm) (int, domain.ResponseLogin) {
	response := domain.ResponseLogin{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	config, err := util.LoadConfig()
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	memberBody, err := m.memberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Code = domain.ErrMemberEmailNotFound.Code
		response.Message = domain.ErrMemberEmailNotFound.Msg
		return http.StatusBadRequest, response
	}

	// Check password
	if memberBody.Password != core_util.HashPassword(loginForm.Password, config.KeyHashPassword) {
		response.Code = domain.ErrLoginHashPassword.Code
		response.Message = domain.ErrLoginHashPassword.Msg
		return http.StatusUnauthorized, response
	}

	// create token
	hashAuth := domain.HashAuth{
		Username: memberBody.Username,
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}

	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return http.StatusOK, response
}

func (m *MemberUtil) GoogleRegister(responseGoogle any) (int, domain.ResponseRegisterMember) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response := domain.ResponseRegisterMember{}
	info, ok := responseGoogle.(domain.ResponseGoogleRegisterCallback)
	if !ok {
		response.Code = domain.ErrInvalidOAuthDecode.Code
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusUnauthorized, response
	}

	config, err := util.LoadConfig()
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return 500, response
	}

	createAt := time.Now()
	updateAt := time.Now()
	password := util.Sha256(info.UserInfo.ID + google)
	memberBody := domain.Member{
		UserID:       core_util.GenUserID(info.UserInfo.Email, password),
		Email:        info.UserInfo.Email,
		Hash:         core_util.GenerateHash(config.KeyHashMember),
		DisplayName:  info.UserInfo.Name,
		ProfileImage: info.UserInfo.Picture,
		GoogleID:     info.UserInfo.ID,
		CreatedAt:    createAt,
		UpdatedAt:    updateAt,
	}
	if err := m.memberRepo.RegisterMember(ctx, memberBody); err != nil {
		response.Code = domain.ErrCreateMemberFail.Code
		response.Message = domain.ErrCreateMemberFail.Msg
		return 400, response
	}

	// create token
	hashAuth := domain.HashAuth{
		Username: memberBody.Username,
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return 500, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return 200, response
}

func (m *MemberUtil) GoogleLogin(responseGoogle any) (int, domain.ResponseLogin) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response := domain.ResponseLogin{}
	info, ok := responseGoogle.(domain.ResponseGoogleLoginCallback)
	if !ok {
		response.Code = domain.ErrInvalidOAuthDecode.Code
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusBadRequest, response
	}

	loginForm := domain.LoginForm{
		Email: info.UserInfo.Email,
	}
	memberBody, err := m.memberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Code = domain.ErrMemberEmailNotFound.Code
		response.Message = domain.ErrMemberEmailNotFound.Msg
		return http.StatusBadRequest, response
	}

	// create token
	hashAuth := domain.HashAuth{
		Username: memberBody.Username,
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryptetion())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.BearerToken = token
	response.Message = "Success"
	return http.StatusOK, response
}

type RedirectUtil struct {
	cache port.Cache
}

func NewRedirectUtil(cache port.Cache) *RedirectUtil {
	redirecUtil := &RedirectUtil{
		cache: cache,
	}
	return redirecUtil
}

func (m *RedirectUtil) GoogleLogin() (int, domain.ResponseRedirectGoogleLogin) {
	response := domain.ResponseRedirectGoogleLogin{}
	googleConfig, err := util.GoogleConfig(login)
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	if googleConfig.State, err = core_util.RandomGoogleState(); err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	ltState := time.Duration(5 * time.Minute)
	if err := m.cache.SetGoogleState(googleConfig.State, ltState); err != nil {
		response.Code = domain.ErrSetGoogleState.Code
		response.Message = domain.ErrSetGoogleState.Msg
		return http.StatusInternalServerError, response
	}

	response.URL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	return http.StatusTemporaryRedirect, response
}

func (m *RedirectUtil) GoogleRegister() (int, domain.ResponseRedirectGoogleRegister) {
	response := domain.ResponseRedirectGoogleRegister{}
	googleConfig, err := util.GoogleConfig(register)
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	if googleConfig.State, err = core_util.RandomGoogleState(); err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	ltState := time.Duration(5 * time.Minute)
	if err := m.cache.SetGoogleState(googleConfig.State, ltState); err != nil {
		response.Code = domain.ErrSetGoogleState.Code
		response.Message = domain.ErrSetGoogleState.Msg
		return http.StatusInternalServerError, response
	}

	response.URL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	return http.StatusTemporaryRedirect, response
}
