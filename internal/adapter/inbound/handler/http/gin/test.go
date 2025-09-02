package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/util"
	"encoding/json"
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

func TestLogin(c *gin.Context) {
	typeLogin, ok := c.Get("type_login")
	fmt.Println(ok, " : ", typeLogin)

	rawResponse, ok := c.Get("response")
	fmt.Println(ok, " : ", rawResponse)

	byteTest, _ := json.Marshal(rawResponse)
	fmt.Println("Raw response JSON:", string(byteTest))

	byteResponse, _ := util.EncryptGOB(rawResponse)
	response := domain.ResponseGoogleLoginCallback{}
	_ = util.DecryptGOB(byteResponse, &response)
	// byteResponse, _ := json.Marshal(rawResponse)
	// json.Unmarshal(byteResponse, &response)

	fmt.Printf("Decrypted response: %#v\n", response)
}
