package gin

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *MainRoute) RouteTest() {
	m.Engine.GET("/test", TestFunc)
}

func TestFunc(c *gin.Context) {
	fmt.Println("Test function called")
	for i := 0; i < 20; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	c.JSON(200, gin.H{"message": "Test function called"})
}
