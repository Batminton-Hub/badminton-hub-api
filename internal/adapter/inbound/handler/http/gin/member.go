package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type MemberController interface {
	RegisterMember(c *gin.Context)
	Login(c *gin.Context)
	GoogleLogin(c *gin.Context)
	GoogleLoginCallback(c *gin.Context)
	GoogleRegister(c *gin.Context)
	GoogleRegisterCallback(c *gin.Context)
}
type MemberControllerImpl struct {
	MemberUtil port.MemberUtil
}

func (m *MemberControllerImpl) RegisterMember(c *gin.Context) {
	registerForm := domain.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	httpStatus, response := m.MemberUtil.RegisterMember(registerForm)

	c.JSON(httpStatus, response)
}

func (m *MemberControllerImpl) Login(c *gin.Context) {
	loginForm := domain.LoginForm{}
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	httpStatus, response := m.MemberUtil.Login(loginForm)

	c.JSON(httpStatus, response)
}

func (m *MemberControllerImpl) GoogleLogin(c *gin.Context) {

	httpStatus, response := m.MemberUtil.GoogleLogin()

	c.Redirect(httpStatus, response.URL)
}

func (m *MemberControllerImpl) GoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MemberUtil.GoogleLoginCallback(state, code)
	if httpStatus != http.StatusOK {
		c.AbortWithStatus(httpStatus)
		return
	}

	c.Set("response", response)
	c.Set("type_login", "google")

	c.Next()
}

func (m *MemberControllerImpl) GoogleRegister(c *gin.Context) {
	googleConfig, err := util.GoogleConfig("REGISTER")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get Google config"})
		return
	}
	url := googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (m *MemberControllerImpl) GoogleRegisterCallback(c *gin.Context) {
	fmt.Println("GoogleRegisterCallback called")
	fmt.Println("Query Parameters:", c.Request.URL.Query())
	prompt := c.Query("prompt")
	fmt.Println("Prompt Parameter:", prompt)

	typeGoogle := c.Query("type")
	fmt.Println("Type Parameter:", typeGoogle)

	c.JSON(200, gin.H{
		"message": "Google register callback received",
		"prompt":  prompt,
		"type":    typeGoogle,
	})
}
