package gin

import "github.com/gin-gonic/gin"

func RespError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
