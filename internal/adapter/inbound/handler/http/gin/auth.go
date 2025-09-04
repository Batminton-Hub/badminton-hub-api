package gin

import (
	"Badminton-Hub/internal/core/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const google = "GOOGLE"

type MiddlewareControllerImpl struct {
	port.MiddlewareUtil
}

func (m *MiddlewareControllerImpl) Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	code, resp := m.MiddlewareUtil.Authenticate(token)
	if code != 200 {
		c.AbortWithStatusJSON(code, resp)
		return
	}

	c.Next()
}

func (m *MiddlewareControllerImpl) GoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MiddlewareUtil.GoogleLoginCallback(state, code)
	if httpStatus != http.StatusOK {
		fmt.Println("GoogleLoginCallback error:", response.Message)
		c.AbortWithStatus(httpStatus)
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
		fmt.Println("GoogleRegisterCallback error:", response.Message)
		c.AbortWithStatus(httpStatus)
		return
	}

	c.Set("response", response)
	c.Set("platform", google)

	c.Next()
}
