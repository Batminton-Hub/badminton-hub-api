package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"fmt"

	"github.com/gin-gonic/gin"
)

type RedirectController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type Redirect struct {
	redirect port.RedirectService
}

func (r *Redirect) Login(c *gin.Context) {
	fmt.Println(getPlatformParam(c))
	redirectInfo := domain.RedirectLoginInfo{
		Platform: getPlatformParam(c),
	}
	httpStatus, response := r.redirect.Login(redirectInfo)
	if response.Resp.Status == domain.ERROR {
		Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, nil)
		return
	}

	RespRedirect(c, httpStatus, response.Resp.Code, response.Resp.Msg, response.URL)
}

func (r *Redirect) Register(c *gin.Context) {
	fmt.Println(getPlatformParam(c))
	redirectInfo := domain.RedirectLoginInfo{
		Platform: getPlatformParam(c),
	}
	httpStatus, response := r.redirect.Register(redirectInfo)
	if response.Resp.Status == domain.ERROR {
		Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, nil)
		return
	}

	RespRedirect(c, httpStatus, response.Resp.Code, response.Resp.Msg, response.URL)
}
