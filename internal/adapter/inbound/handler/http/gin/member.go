package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberControllerImpl struct {
	MemberUtil port.MemberService
}

func (m *MemberControllerImpl) RegisterMember(c *gin.Context) {
	var httpStatus int
	var response domain.ResponseRegisterMember
	var registerForm domain.RegisterForm
	platform := c.GetString("platform")

	switch platform {
	case "GOOGLE":
		responseGoogle, ok := c.Get("response")
		if !ok {
			RespError(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg)
			return
		}
		httpStatus, response = m.MemberUtil.GoogleRegister(responseGoogle)
	default:
		if err := c.ShouldBind(&registerForm); err != nil {
			RespError(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg)
			return
		}
		httpStatus, response = m.MemberUtil.RegisterMember(registerForm)
	}

	if httpStatus != http.StatusOK {
		RespError(c, httpStatus, response.Code, response.Message)
		return
	}
	RespSuccess(c, httpStatus, response)
}

func (m *MemberControllerImpl) Login(c *gin.Context) {
	var httpStatus int
	var response domain.ResponseLogin
	var loginForm domain.LoginForm
	platform := c.GetString("platform")

	switch platform {
	case "GOOGLE":
		responseGoogle, ok := c.Get("response")
		if !ok {
			RespError(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg)
			return
		}
		httpStatus, response = m.MemberUtil.GoogleLogin(responseGoogle)
	default:
		if err := c.ShouldBind(&loginForm); err != nil {
			RespError(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg)
			return
		}
		httpStatus, response = m.MemberUtil.Login(loginForm)
	}

	if httpStatus != http.StatusOK {
		RespError(c, httpStatus, response.Code, response.Message)
		return
	}
	RespSuccess(c, httpStatus, response)
}

func (m *MemberControllerImpl) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	fmt.Println("userID", userID)
	httpStatus, response := m.MemberUtil.GetProfile(userID)
	if httpStatus != http.StatusOK {
		RespError(c, httpStatus, response.Code, response.Message)
		return
	}

	RespSuccess(c, httpStatus, response)
}

func (m *MemberControllerImpl) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	fmt.Println("userID", userID)
	var request domain.RequestUpdateProfile
	if err := c.ShouldBind(&request); err != nil {
		RespError(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg)
		return
	}

	httpStatus, response := m.MemberUtil.UpdateProfile(userID, request)
	if httpStatus != http.StatusOK {
		RespError(c, httpStatus, response.Code, response.Message)
		return
	}

	RespSuccess(c, httpStatus, response)
}

type RedirectControllerImpl struct {
	RedirectUtil port.RedirectUtil
}

func (m *RedirectControllerImpl) GoogleLogin(c *gin.Context) {
	httpStatus, response := m.RedirectUtil.GoogleLogin()
	if response.Code != 0 {
		RespError(c, httpStatus, response.Code, response.Message)
		return
	}
	c.Redirect(httpStatus, response.URL)
}

func (m *RedirectControllerImpl) GoogleRegister(c *gin.Context) {
	httpStatus, response := m.RedirectUtil.GoogleRegister()
	if response.Code != 0 {
		RespError(c, httpStatus, response.Code, response.Message)
		return
	}

	c.Redirect(httpStatus, response.URL)
}
