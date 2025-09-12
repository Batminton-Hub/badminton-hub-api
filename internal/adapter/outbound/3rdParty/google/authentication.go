package google

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"encoding/json"
	"net/http"
	"time"
)

type GoogleMiddleware struct {
	cache port.CacheUtil
}

func NewGoogleMiddleware(
	cache port.CacheUtil,
) *GoogleMiddleware {
	return &GoogleMiddleware{
		cache: cache,
	}
}

func (g *GoogleMiddleware) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	ctx, cancel := util.InitConText(10 * time.Second)
	defer cancel()

	state := info.State
	code := info.Code
	platformData := GoogleMemberInfo{}
	response := domain.RespAuth{}

	googleConfig, err := GoogleConfig(info.Action)
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

	var userInfo GoogleUserInfo
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
