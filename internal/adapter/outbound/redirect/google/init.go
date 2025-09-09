package google

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/internal/core_util"
	"Badminton-Hub/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type GoogleRedirect struct {
	cache port.CacheUtil
}

func NewGoogleRedirect(
	cache port.CacheUtil,
) *GoogleRedirect {
	return &GoogleRedirect{
		cache: cache,
	}
}

func (g *GoogleRedirect) Login(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	var response domain.RespRedirect
	var redirectURL string
	googleConfig, err := util.GoogleConfig(domain.LOGIN)
	if err != nil {
		response.Resp = domain.ErrLoadConfig
		return response.Resp.HttpStatus, response
	}

	if googleConfig.State, err = core_util.RandomGoogleState(); err != nil {
		fmt.Println("invalid oauth state[1] :", err)
		response.Resp = domain.ErrInvalidOAuthState
		return response.Resp.HttpStatus, response
	}

	ltState := time.Duration(5 * time.Minute)
	if err := g.cache.SetGoogleState(ctx, googleConfig.State, ltState); err != nil {
		response.Resp = domain.ErrSetGoogleState
		return response.Resp.HttpStatus, response
	}

	redirectURL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	response.URL = redirectURL
	response.Resp = domain.RedirectSuccess
	return response.Resp.HttpStatus, response
}

func (g *GoogleRedirect) Register(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	var response domain.RespRedirect
	var redirectURL string
	googleConfig, err := util.GoogleConfig(domain.REGISTER)
	if err != nil {
		response.Resp = domain.ErrLoadConfig
		return response.Resp.HttpStatus, response
	}

	if googleConfig.State, err = core_util.RandomGoogleState(); err != nil {
		fmt.Println("invalid oauth state[2] :", err)
		response.Resp = domain.ErrInvalidOAuthState
		return response.Resp.HttpStatus, response
	}

	ltState := time.Duration(5 * time.Minute)
	if err := g.cache.SetGoogleState(ctx, googleConfig.State, ltState); err != nil {
		response.Resp = domain.ErrSetGoogleState
		return response.Resp.HttpStatus, response
	}

	redirectURL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	response.URL = redirectURL
	response.Resp = domain.RedirectSuccess
	return response.Resp.HttpStatus, response
}

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
		fmt.Println("invalid oauth exchange[1] :", err)
		response.Resp = domain.ErrInvalidOAuthExchange
		return http.StatusUnauthorized, response
	}

	// ใช้ token เรียก API ของ Google เพื่อดึงข้อมูลผู้ใช้
	client := googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
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
