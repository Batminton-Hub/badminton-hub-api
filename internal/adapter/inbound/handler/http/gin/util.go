package gin

import (
	"Badminton-Hub/internal/core/domain"

	"github.com/gin-gonic/gin"
)

func RespSuccess(c *gin.Context, httpStaus int, response any) {
	c.JSON(httpStaus, response)
}

func RespError(c *gin.Context, httpStaus int, code int, message string) {
	response := domain.ResponseError{
		Code:    code,
		Message: message,
	}
	c.JSON(httpStaus, response)
}

func RespAuth(c *gin.Context, httpStaus, code int, message string) {
	response := domain.ResponseError{
		Code:    code,
		Message: message,
	}
	c.AbortWithStatusJSON(httpStaus, response)
}
