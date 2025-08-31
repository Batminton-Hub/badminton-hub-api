package util

import (
	"Badminton-Hub/internal/core/domain"
	"fmt"
	"time"
)

var encryption = NewJWTEncryption()

type AuthBody struct {
	Exp        int64             `json:"exp"`
	Data       domain.AuthMember `json:"data"`
	Permission map[string]int    `json:"permission"`
}

func GenBearerToken(member domain.Member) (string, error) {
	lt := time.Duration(5 * time.Minute)
	exp := time.Now().Add(lt).Unix()
	createAt := time.Now().UTC()
	rawHash := member.Email + member.Username + fmt.Sprint(createAt) + "test"
	authBody := AuthBody{
		Exp: exp,
		Data: domain.AuthMember{
			CreatedAt: createAt,
			HashAuth:  HashAuth(rawHash),
		},
	}

	encryptedMember, err := encryption.Encrypte(authBody, "your-encryption-key-here", lt)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt member: %w", err)
	}

	token := encryptedMember

	return token, nil
}

func ValidateBearerToken(token string) (AuthBody, error) {
	authBody := AuthBody{}

	err := encryption.Decrypte(token, "your-encryption-key-here", &authBody)
	if err != nil {
		return authBody, err
	}

	if authBody.Exp < time.Now().Unix() {
		return authBody, fmt.Errorf("token has expired")
	}

	return authBody, nil
}
