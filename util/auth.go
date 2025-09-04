package util

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"fmt"
	"math/rand"
	"time"
)

type AuthBody struct {
	Exp        int64             `json:"exp"`
	Data       domain.AuthMember `json:"data"`
	Permission map[string]int    `json:"permission"`
}

func GenBearerToken(hashBody domain.HashAuth, encryption port.Encryption) (string, error) {
	var token string
	config, err := LoadConfig()
	if err != nil {
		return token, fmt.Errorf("failed to load config: %w", err)
	}

	lt := time.Duration(5 * time.Minute)
	exp := time.Now().Add(lt).Unix()
	createAt := time.Now().UTC()
	byteHash, err := EncryptGOB(hashBody)
	if err != nil {
		return token, fmt.Errorf("failed to encrypt hash body: %w", err)
	}
	rawHash := string(byteHash)
	authBody := AuthBody{
		Exp: exp,
		Data: domain.AuthMember{
			CreatedAt: createAt,
			HashAuth:  HashAuth(rawHash, config.KeyHashAuth),
		},
	}

	encryptedMember, err := encryption.Encrypte(authBody, "your-encryption-key-here", lt)
	if err != nil {
		return token, fmt.Errorf("failed to encrypt member: %w", err)
	}

	token = encryptedMember

	return token, nil
}

func ValidateBearerToken(encryption port.Encryption, token string) (AuthBody, error) {
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

func RandomGoogleState() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if config.Mode == "DEVERLOP" {
		return "0123456789ABCDEF", nil
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	state := make([]byte, 32)
	for i := range state {
		state[i] = charset[rand.Intn(len(charset))]
	}
	return string(state), nil
}
