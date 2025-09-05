package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	core_util "Badminton-Hub/internal/util"
	"Badminton-Hub/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MiddlewareUtil struct {
	encryption port.Encryption
	cache      port.Cache
}

func (m *MiddlewareUtil) Encryption() port.Encryption { return m.encryption }

func NewMiddlewareUtil(encryption port.Encryption, cache port.Cache) *MiddlewareUtil {
	return &MiddlewareUtil{
		encryption: encryption,
		cache:      cache,
	}
}

func (m *MiddlewareUtil) Authenticate(token string) (int, domain.AuthResponse) {
	response := domain.AuthResponse{}
	config := util.LoadConfig()

	// Remove "Bearer " prefix
	token = token[len("Bearer "):]

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := core_util.ValidateBearerToken(token, m.encryption)
	if err != nil {
		fmt.Println("Authenticate err : ", err)
		response.Code = domain.ErrValidateToken.Code
		response.Message = domain.ErrValidateToken.Msg
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
	hashauth := core_util.HashAuth(rawHash, config.KeyHashAuth)
	if authBody.Data.HashAuth != hashauth {
		response.Code = domain.ErrValidateHashAuth.Code
		response.Message = domain.ErrValidateHashAuth.Msg
		return http.StatusUnauthorized, response
	}

	// ตรวจสอบว่า token ยังไม่หมดอายุ
	if authBody.Exp < time.Now().Unix() {
		response.Code = domain.ErrTokenExpired.Code
		response.Message = domain.ErrTokenExpired.Msg
		return http.StatusUnauthorized, response
	}

	response.Code = domain.AuthSuccess.Code
	response.Message = domain.AuthSuccess.Msg
	return http.StatusOK, response
}

func (m *MiddlewareUtil) GoogleLoginCallback(state, code string) (int, domain.ResponseGoogleLoginCallback) {
	fmt.Println("GoogleLoginCallback Start")
	ctx, cancel := util.InitConText(10 * time.Second)
	defer cancel()

	response := domain.ResponseGoogleLoginCallback{}
	googleConfig, err := util.GoogleConfig(LOGIN)
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	googleOAuthConfig := googleConfig.Config

	// check and delete state
	if ok, err := m.cache.GetGoogleState(ctx, state); err != nil || !ok {
		response.Code = domain.ErrInvalidOAuthState.Code
		response.Message = domain.ErrInvalidOAuthState.Msg
		return http.StatusUnauthorized, response
	}
	m.cache.DelGoogleState(ctx, state)

	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		response.Code = domain.ErrInvalidOAuthExchange.Code
		response.Message = domain.ErrInvalidOAuthExchange.Msg
		return http.StatusUnauthorized, response
	}

	// ใช้ token เรียก API ของ Google เพื่อดึงข้อมูลผู้ใช้
	client := googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		response.Code = domain.ErrInvalidOAuthClient.Code
		response.Message = domain.ErrInvalidOAuthClient.Msg
		return http.StatusUnauthorized, response
	}
	defer resp.Body.Close()

	var userInfo domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response.Code = domain.ErrInvalidOAuthDecode.Code
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusUnauthorized, response
	}

	response.UserInfo = userInfo
	response.AccessToken = token.AccessToken
	response.RefreshToken = token.RefreshToken
	response.Code = domain.AuthSuccess.Code
	response.Message = domain.AuthSuccess.Msg
	return http.StatusOK, response
}

func (m *MiddlewareUtil) GoogleRegisterCallback(state, code string) (int, domain.ResponseGoogleRegisterCallback) {
	ctx, cancel := util.InitConText(10 * time.Second)
	defer cancel()

	response := domain.ResponseGoogleRegisterCallback{}
	googleConfig, err := util.GoogleConfig(REGISTER)
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	googleOAuthConfig := googleConfig.Config

	// check and delete state
	if ok, err := m.cache.GetGoogleState(ctx, state); err != nil || !ok {
		response.Code = domain.ErrInvalidOAuthState.Code
		response.Message = domain.ErrInvalidOAuthState.Msg
		return http.StatusUnauthorized, response
	}
	m.cache.DelGoogleState(ctx, state)

	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		response.Code = domain.ErrInvalidOAuthExchange.Code
		response.Message = domain.ErrInvalidOAuthExchange.Msg
		return http.StatusUnauthorized, response
	}

	// ใช้ token เรียก API ของ Google เพื่อดึงข้อมูลผู้ใช้
	client := googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		response.Code = domain.ErrInvalidOAuthClient.Code
		response.Message = domain.ErrInvalidOAuthClient.Msg
		return http.StatusUnauthorized, response
	}
	defer resp.Body.Close()

	var userInfo domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response.Code = domain.ErrInvalidOAuthDecode.Code
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusUnauthorized, response
	}

	response.UserInfo = userInfo
	response.AccessToken = token.AccessToken
	response.RefreshToken = token.RefreshToken
	response.Code = domain.AuthSuccess.Code
	response.Message = domain.AuthSuccess.Msg
	return http.StatusOK, response
}
