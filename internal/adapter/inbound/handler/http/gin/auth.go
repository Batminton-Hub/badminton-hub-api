package gin

import (
	"Badminton-Hub/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

const google = "GOOGLE"

type MiddlewareControllerImpl struct {
	port.MiddlewareService
}

func (m *MiddlewareControllerImpl) Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	httpStatus, response := m.MiddlewareService.Authenticate(token)
	if httpStatus != http.StatusOK {
		RespAuth(c, httpStatus, response.Code, response.Message)
		return
	}

	c.Set("user_id", response.AuthBody.Data.UserID)

	c.Next()
}

func (m *MiddlewareControllerImpl) GoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MiddlewareService.GoogleLoginCallback(state, code)
	if httpStatus != http.StatusOK {
		RespAuth(c, httpStatus, response.Code, response.Message)
		return
	}

	c.Set("response", response)
	c.Set("platform", google)

	c.Next()
}

func (m *MiddlewareControllerImpl) GoogleRegisterCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MiddlewareService.GoogleRegisterCallback(state, code)
	if httpStatus != http.StatusOK {
		RespAuth(c, httpStatus, response.Code, response.Message)
		return
	}

	c.Set("response", response)
	c.Set("platform", google)

	c.Next()
}
