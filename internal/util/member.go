package core_util

import (
	"Badminton-Hub/util"
	"fmt"
	"time"
)

func GenerateHash(key string) string {
	timeNow := time.Now().UnixMilli()
	data := fmt.Sprint(key, timeNow)
	hash := util.Sha256(data)
	return hash
}
func HashPassword(password, key string) string {
	data := fmt.Sprint(password, key)
	newPassword := util.Sha256(data)
	return newPassword
}

func GenUserID(email, password string) string {
	return util.Sha256(email + password)
}
