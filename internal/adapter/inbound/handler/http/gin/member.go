package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberControllerImpl struct {
	MemberUtil port.MemberUtil
}

func (m *MemberControllerImpl) RegisterMember(c *gin.Context) {
	var httpStatus int
	var response domain.ResponseRegisterMember
	var registerForm domain.RegisterForm
	platform := c.GetString("platform")

	fmt.Println("platform", platform)
	switch platform {
	case "GOOGLE":
		responseGoogle, ok := c.Get("response")
		if !ok {
			Resp(c, http.StatusBadRequest, "Invalid input")
			return
		}
		httpStatus, response = m.MemberUtil.GoogleRegister(responseGoogle)
	default:
		if err := c.ShouldBind(&registerForm); err != nil {
			Resp(c, http.StatusBadRequest, "Invalid input")
			return
		}
		httpStatus, response = m.MemberUtil.RegisterMember(registerForm)
	}

	c.JSON(httpStatus, response)
}

func (m *MemberControllerImpl) Login(c *gin.Context) {
	var httpStatus int
	var response domain.ResponseLogin
	var loginForm domain.LoginForm
	platform := c.GetString("platform")

	fmt.Println("platform", platform)
	switch platform {
	case "GOOGLE":
		responseGoogle, ok := c.Get("response")
		if !ok {
			Resp(c, http.StatusBadRequest, "Invalid input")
			return
		}
		httpStatus, response = m.MemberUtil.GoogleLogin(responseGoogle)
	default:
		if err := c.ShouldBind(&loginForm); err != nil {
			Resp(c, http.StatusBadRequest, "Invalid input")
			return
		}
		httpStatus, response = m.MemberUtil.Login(loginForm)
	}

	c.JSON(httpStatus, response)
}

type RedirectControllerImpl struct {
	RedirectUtil port.RedirectUtil
}

func (m *RedirectControllerImpl) GoogleLogin(c *gin.Context) {
	httpStatus, response := m.RedirectUtil.GoogleLogin()
	c.Redirect(httpStatus, response.URL)
}

func (m *RedirectControllerImpl) GoogleRegister(c *gin.Context) {
	httpStatus, response := m.RedirectUtil.GoogleRegister()
	if httpStatus != http.StatusOK {
		Resp(c, httpStatus, response)
	}
	c.Redirect(httpStatus, response.URL)
}

type ProfileControllerImpl struct {
	MemberUtil port.MemberUtil
}

func (m *ProfileControllerImpl) GetProfile(c *gin.Context) {
	// token := c.GetHeader("Authorization")
	// code, resp := m.MemberUtil.GetProfile(token)
	// if code != 200 {
	// 	c.AbortWithStatusJSON(code, resp)
	// 	return
	// }
}
