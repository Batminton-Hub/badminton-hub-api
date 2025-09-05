package core_util

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"fmt"
	"time"
)

func HashAuth(rawHash, key string) string {
	data := fmt.Sprint(rawHash + key)
	hashAuth := util.Sha256(data)
	return hashAuth
}

func GenBearerToken(hashBody domain.HashAuth, encryption port.Encryption) (string, error) {
	var token string
	config := util.LoadConfig()

	lt := time.Duration(5 * time.Minute)
	exp := time.Now().Add(lt).Unix()
	createAt := time.Now().UTC()
	byteHash, err := util.EncryptGOB(hashBody)
	if err != nil {
		return token, fmt.Errorf("failed to encrypt hash body: %w", err)
	}
	rawHash := string(byteHash)
	authBody := domain.AuthBody{
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

func ValidateBearerToken(encryption port.Encryption, token string) (domain.AuthBody, error) {
	authBody := domain.AuthBody{}

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
	config := util.LoadConfig()

	if config.Mode == "DEVERLOP" {
		return config.DefaultGoogleState, nil
	}

	state := util.RandomString(32, true, true, false)
	return string(state), nil
}
