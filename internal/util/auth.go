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

func GenBearerToken(hashBody domain.HashAuth, encryption port.EncryptionUtil) (string, error) {
	var token string
	config := util.LoadConfig()

	lt := config.BearerTokenExp
	createAt := time.Now().UTC()
	exp := time.Now().Add(lt).Unix()

	byteHash, err := util.EncryptGOB(hashBody)
	if err != nil {
		return token, fmt.Errorf("failed to encrypt hash body: %w", err)
	}
	rawHash := string(byteHash)
	authBody := domain.AuthBody{
		CreateAt: createAt,
		Exp:      exp,
		Data: domain.AuthMember{
			UserID:    hashBody.UserID,
			CreatedAt: hashBody.CreateAt,
			HashAuth:  HashAuth(rawHash, config.KeyHashAuth),
		},
	}

	encryptedMember, err := encryption.Encrypte(authBody, config.KeyBearerToken, lt)
	if err != nil {
		return token, fmt.Errorf("failed to encrypt member: %w", err)
	}

	token = encryptedMember

	return token, nil
}

func ValidateBearerToken(token string, encryption port.EncryptionUtil) (domain.AuthBody, error) {
	config := util.LoadConfig()
	authBody := domain.AuthBody{}
	err := encryption.Decrypte(token, config.KeyBearerToken, &authBody)
	if err != nil {
		return authBody, err
	}

	if authBody.Exp < time.Now().Unix() {
		fmt.Println("authBody.Exp < time.Now().Unix() : ", authBody.Exp, " : ", time.Now().Unix())
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
