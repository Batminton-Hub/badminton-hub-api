package google

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"encoding/json"
	"net/http"
	"time"
)

type GoogleCallback struct {
	cache port.CacheUtil
}

func NewGoogleCallback(
	cache port.CacheUtil,
) *GoogleCallback {
	return &GoogleCallback{
		cache: cache,
	}
}

func (g *GoogleCallback) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	ctx, cancel := util.InitConText(10 * time.Second)
	defer cancel()

	state := info.State
	code := info.Code
	platformData := domain.ResponseGoogleLoginCallback{}
	response := domain.RespAuth{}

	googleConfig, err := util.GoogleConfig(info.Action)
	if err != nil {
		response.Resp = domain.ErrLoadConfig
		return http.StatusInternalServerError, response
	}

	googleOAuthConfig := googleConfig.Config

	// check and delete state
	if ok, err := g.cache.GetGoogleState(ctx, state); err != nil || !ok {
		response.Resp = domain.ErrInvalidOAuthState
		return http.StatusUnauthorized, response
	}
	g.cache.DelGoogleState(ctx, state)

	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		response.Resp = domain.ErrInvalidOAuthExchange
		return http.StatusUnauthorized, response
	}

	// ใช้ token เรียก API ของ Google เพื่อดึงข้อมูลผู้ใช้
	client := googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get(domain.GoogleUserInfoURL)
	if err != nil {
		response.Resp = domain.ErrInvalidOAuthClient
		return http.StatusUnauthorized, response
	}
	defer resp.Body.Close()

	var userInfo domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response.Resp = domain.ErrInvalidOAuthDecode
		return http.StatusUnauthorized, response
	}

	platformData.UserInfo = userInfo
	platformData.AccessToken = token.AccessToken
	platformData.RefreshToken = token.RefreshToken
	response.PlatformData = platformData
	response.Resp = domain.AuthSuccess
	return response.Resp.HttpStatus, response
}
