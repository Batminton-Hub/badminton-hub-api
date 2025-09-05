package gin

import (
	"Badminton-Hub/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

const google = "GOOGLE"

type MiddlewareControllerImpl struct {
	port.MiddlewareUtil
}

func (m *MiddlewareControllerImpl) Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	httpCode, response := m.MiddlewareUtil.Authenticate(token)
	if httpCode != http.StatusOK {
		RespAuth(c, httpCode, response)
		return
	}

	c.Next()
}

func (m *MiddlewareControllerImpl) GoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MiddlewareUtil.GoogleLoginCallback(state, code)
	if httpStatus != http.StatusOK {
		RespAuth(c, httpStatus, response)
		return
	}

	c.Set("response", response)
	c.Set("platform", google)

	c.Next()
}

func (m *MiddlewareControllerImpl) GoogleRegisterCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MiddlewareUtil.GoogleRegisterCallback(state, code)
	if httpStatus != http.StatusOK {
		RespAuth(c, httpStatus, response)
		return
	}

	c.Set("response", response)
	c.Set("platform", google)

	c.Next()
}
