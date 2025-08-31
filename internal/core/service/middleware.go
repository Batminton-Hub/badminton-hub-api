package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/util"
	"time"
)

type MiddlewareUtil struct{}

func NewMiddlewareUtil() *MiddlewareUtil {
	return &MiddlewareUtil{}
}

func (m *MiddlewareUtil) Authenticate(token string) (int, domain.AuthResponse) {
	response := domain.AuthResponse{}
	// Remove "Bearer " prefix
	token = token[len("Bearer "):]

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := util.ValidateBearerToken(token)
	if err != nil {
		response.Code = domain.ErrValidateToken.Code
		response.Err = domain.ErrValidateToken.Err
		return 401, response
	}

	// ตรวจสอบความถูกต้องของ token
	byteAuth, err := util.EncryptGOB(authBody)
	if err != nil {
		return 401, response
	}
	rawHash := string(byteAuth)
	hashauth := util.HashAuth(rawHash)
	if authBody.Data.HashAuth != hashauth {
		response.Code = domain.ErrValidateHashAuth.Code
		response.Err = domain.ErrValidateHashAuth.Err
		return 401, response
	}

	// ตรวจสอบว่า token ยังไม่หมดอายุ
	if authBody.Exp < time.Now().Unix() {
		response.Code = domain.ErrTokenExpired.Code
		response.Err = domain.ErrTokenExpired.Err
		return 401, response
	}

	return 200, response
}
