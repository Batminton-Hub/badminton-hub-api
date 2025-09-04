package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	core_util "Badminton-Hub/internal/util"
	"Badminton-Hub/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MiddlewareUtil struct {
	encryption port.Encryption
	cache      port.Cache
}

func (m *MiddlewareUtil) Encryptetion() port.Encryption { return m.encryption }

func NewMiddlewareUtil(encryption port.Encryption, cache port.Cache) *MiddlewareUtil {
	return &MiddlewareUtil{
		encryption: encryption,
		cache:      cache,
	}
}

func (m *MiddlewareUtil) Authenticate(token string) (int, domain.AuthResponse) {
	response := domain.AuthResponse{}
	config, err := util.LoadConfig()
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Err = domain.ErrLoadConfig.Err
		return 500, response
	}

	// Remove "Bearer " prefix
	token = token[len("Bearer "):]

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := core_util.ValidateBearerToken(m.encryption, token)
	if err != nil {
		response.Code = domain.ErrValidateToken.Code
		response.Err = domain.ErrValidateToken.Err
		return 401, response
	}

	// ตรวจสอบความถูกต้องของ token
	byteAuth, err := util.EncryptGOB(authBody)
	if err != nil {
		return 401, response
	}
	rawHash := string(byteAuth)
	hashauth := core_util.HashAuth(rawHash, config.KeyHashAuth)
	if authBody.Data.HashAuth != hashauth {
		response.Code = domain.ErrValidateHashAuth.Code
		response.Err = domain.ErrValidateHashAuth.Err
		return 401, response
	}

	// ตรวจสอบว่า token ยังไม่หมดอายุ
	if authBody.Exp < time.Now().Unix() {
		response.Code = domain.ErrTokenExpired.Code
		response.Err = domain.ErrTokenExpired.Err
		return 401, response
	}

	return 200, response
}

func (m *MiddlewareUtil) GoogleLoginCallback(state, code string) (int, domain.ResponseGoogleLoginCallback) {
	fmt.Println("GoogleLoginCallback Start")
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

	// check state
	if ok, err := m.cache.GetGoogleState(state); err != nil || !ok {
		response.ErrorCode = domain.ErrInvalidOAuthState.Code
		response.Error = domain.ErrInvalidOAuthState.Err
		response.Message = domain.ErrInvalidOAuthState.Msg
		return http.StatusUnauthorized, response
	}
	// m.cache.DelGoogleState(state)

	fmt.Println("Code :", code)
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

	// var userInfo map[string]interface{}
	var userInfo domain.GoogleUserInfo
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
