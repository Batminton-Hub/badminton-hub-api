package gin

import (
	"Badminton-Hub/internal/core/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	Authorization = "Authorization"
	PLATFORM_DATA = "platform_data"
	PLATFORM      = "platform"
	STATE         = "state"
	CODE          = "code"
)

func RespAuth(c *gin.Context, httpStaus, code int, message string, bearerToken string) {
	response := RespAuthBody{
		Code:        code,
		Message:     message,
		BearerToken: bearerToken,
	}
	c.JSON(httpStaus, response)
}

func Resp(c *gin.Context, httpStaus, code int, message string, data any) {
	response := RespBody{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(httpStaus, response)
}

func getPlatform(c *gin.Context) string {
	platform := strings.ToUpper(c.GetString(PLATFORM))
	if platform == "" {
		platform = domain.NORMAL
	}
	return platform
}

func getBearerToken(c *gin.Context) string {
	bearerToken := c.GetHeader(Authorization)
	return bearerToken
}

func getState(c *gin.Context) string {
	state := c.Query(STATE)
	return state
}

func getCode(c *gin.Context) string {
	code := c.Query(CODE)
	return code
}

func getPlatformData(c *gin.Context) any {
	platformData, ok := c.Get(PLATFORM_DATA)
	if !ok {
		return nil
	}
	return platformData
}

func getPlatformParam(c *gin.Context) string {
	return strings.ToUpper(c.Param(PLATFORM))
}

type RespBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RespAuthBody struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	BearerToken string `json:"bearer_token,omitempty"`
}
