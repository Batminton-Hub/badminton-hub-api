package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type MemberUtil struct {
	memberRepo     port.MemberRepo
	middlewareUtil port.MiddlewareUtil
	cache          port.Cache
}

func NewMemberUtil(memberRepo port.MemberRepo, middlewareUtil port.MiddlewareUtil, cache port.Cache) *MemberUtil {
	return &MemberUtil{
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
		cache:          cache,
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

func (m *MemberUtil) GoogleLogin() (int, domain.ResponseGoogleLogin) {
	response := domain.ResponseGoogleLogin{}
	googleConfig, err := util.GoogleConfig("LOGIN")
	if err != nil {
		response.ErrorCode = domain.ErrLoadConfig.Code
		response.Error = domain.ErrLoadConfig.Err
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	googleConfig.State = util.RandomGoogleState()

	response.URL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
	)
	return http.StatusTemporaryRedirect, response
}

func (m *MemberUtil) GoogleLoginCallback(state, code string) (int, domain.ResponseGoogleLoginCallback) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response := domain.ResponseGoogleLoginCallback{}
	googleConfig, err := util.GoogleConfig("LOGIN")
	if err != nil {
		response.ErrorCode = domain.ErrLoadConfig.Code
		response.Error = domain.ErrLoadConfig.Err
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	googleOAuthConfig := googleConfig.Config
	oauthStateString := googleConfig.State

	if state != oauthStateString {
		response.ErrorCode = domain.ErrInvalidOAuthState.Code
		response.Error = domain.ErrInvalidOAuthState.Err
		response.Message = domain.ErrInvalidOAuthState.Msg
		return http.StatusUnauthorized, response
	}

	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		response.ErrorCode = domain.ErrInvalidOAuthExchange.Code
		response.Error = domain.ErrInvalidOAuthExchange.Err
		response.Message = domain.ErrInvalidOAuthExchange.Msg
		return http.StatusUnauthorized, response
	}

	// ใช้ token เรียก API ของ Google เพื่อดึงข้อมูลผู้ใช้
	client := googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		response.ErrorCode = domain.ErrInvalidOAuthClient.Code
		response.Error = domain.ErrInvalidOAuthClient.Err
		response.Message = domain.ErrInvalidOAuthClient.Msg
		return http.StatusUnauthorized, response
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response.ErrorCode = domain.ErrInvalidOAuthDecode.Code
		response.Error = domain.ErrInvalidOAuthDecode.Err
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusUnauthorized, response
	}

	response.UserInfo = userInfo
	response.AccessToken = token.AccessToken
	response.RefreshToken = token.RefreshToken
	return http.StatusOK, response
}

func (m *MemberUtil) GoogleRegister() (int, domain.ResponseGoogleRegister) {
	return 200, domain.ResponseGoogleRegister{}
}

func (m *MemberUtil) GoogleRegisterCallback(state, code string) (int, domain.ResponseGoogleRegisterCallback) {
	return 200, domain.ResponseGoogleRegisterCallback{}
}
