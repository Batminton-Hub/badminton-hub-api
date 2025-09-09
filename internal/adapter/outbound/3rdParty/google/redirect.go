package google

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/internal/core_util"
	"Badminton-Hub/util"
	"time"

	"golang.org/x/oauth2"
)

const testVar = "testVal"

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
