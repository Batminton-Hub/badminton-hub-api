package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"time"
)

type MiddlewareUtil struct {
	Encryption port.Encryption
}

func NewMiddlewareUtil(encryption port.Encryption) *MiddlewareUtil {
	return &MiddlewareUtil{
		Encryption: encryption,
	}
}

func (m *MiddlewareUtil) Authenticate(token string) (int, domain.AuthResponse) {
	response := domain.AuthResponse{}
	config, err := util.LoadConfig()
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Err = domain.ErrLoadConfig.Err
		return 500, response
	}

	// Remove "Bearer " prefix
	token = token[len("Bearer "):]

	// ถอด authentication token ที่ส่งมาจาก client
	authBody, err := util.ValidateBearerToken(m.Encryption, token)
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
	hashauth := util.HashAuth(rawHash, config.KeyHashAuth)
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

func (m *MiddlewareUtil) Encryptetion() port.Encryption {
	return m.Encryption
}
