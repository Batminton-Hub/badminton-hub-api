package gin

import (
	"Badminton-Hub/internal/core/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

func RespAuth(c *gin.Context, httpStaus, code int, message string, bearerToken string) {
	response := RespAuthBody{
		Code:        code,
		Message:     message,
		BearerToken: bearerToken,
	}
	c.JSON(httpStaus, response)
}

func RespMiddleWare(c *gin.Context, httpStatus, code int, message string) {
	response := RespBody{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	c.AbortWithStatusJSON(httpStatus, response)
}

func Resp(c *gin.Context, httpStatus, code int, message string, data any) {
	response := RespBody{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(httpStatus, response)
}

func RespRedirect(c *gin.Context, httpStatus, code int, message string, url string) {
	api := c.Query(domain.API)
	if api != "" {
		Resp(c, httpStatus, code, message, url)
		return
	}
	c.Redirect(httpStatus, url)
}

func getPlatform(c *gin.Context) string {
	platform := strings.ToUpper(c.GetString(domain.Platform))
	if platform == "" {
		platform = domain.NORMAL
	}
	return platform
}

func getBearerToken(c *gin.Context) string {
	bearerToken := c.GetHeader(domain.Authorization)
	return bearerToken
}

func getState(c *gin.Context) string {
	state := c.Query(domain.State)
	return state
}

func getCode(c *gin.Context) string {
	code := c.Query(domain.Code)
	return code
}

func getPlatformData(c *gin.Context) any {
	platformData, ok := c.Get(domain.PlatformData)
	if !ok {
		return nil
	}
	return platformData
}

func getPlatformParam(c *gin.Context) string {
	platform := strings.ToUpper(c.Param(domain.Platform))
	if platform == "" {
		platform = domain.NORMAL
	}
	return platform
}

func getAction(c *gin.Context) string {
	path := strings.ToUpper(c.FullPath())
	switch {
	case strings.Contains(path, domain.LOGIN):
		return domain.LOGIN
	case strings.Contains(path, domain.REGISTER):
		return domain.REGISTER
	default:
		return ""
	}
}

func getUserID(c *gin.Context) string {
	return c.GetString(domain.UserID)
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
