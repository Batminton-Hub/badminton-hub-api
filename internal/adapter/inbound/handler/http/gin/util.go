package gin

import (
	"github.com/gin-gonic/gin"
)

func Resp(c *gin.Context, httpStaus int, response any) {
	c.JSON(httpStaus, response)
}

func RespAuth(c *gin.Context, httpStaus int, response any) {
	c.AbortWithStatusJSON(httpStaus, response)
}
