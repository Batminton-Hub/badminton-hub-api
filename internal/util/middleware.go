package core_util

import (
	"Badminton-Hub/util"
	"errors"
	"fmt"
)

type MiddlewareUtil struct {
}

func NewMiddlewareUtil() *MiddlewareUtil {
	return &MiddlewareUtil{}
}

func (m *MiddlewareUtil) Authenticate(token string) error {

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := util.ValidateBearerToken(token)
	if err != nil {
		return err
	}

	// ตรวจสอบความถูกต้องของ token
	rawHash := authBody.Data.Email + authBody.Data.Username + fmt.Sprint(authBody.Data.CreatedAt) + "test"
	hashauth := util.HashAuth(rawHash)
	if authBody.Data.HashAuth != hashauth {
		return errors.New("invalid token")
	}

	// // ตรวจสอบว่า token ยังไม่หมดอายุ
	// if authBody.Exp < time.Now().Unix() {
	// 	return errors.New("token has expired")
	// }

	// ถ้า token ถูกต้องและยังไม่หมดอายุ ให้ดำเนินการต่อ
	return nil
}
