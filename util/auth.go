package util

import (
	"Badminton-Hub/internal/core/domain"
	"encoding/json"
	"fmt"
	"time"
)

type AuthBody struct {
	Exp  int64             `json:"exp"`
	Data domain.AuthMember `json:"data"`
	// Data interface{} `json:"data"`
}

func GenBearerToken(member domain.Member) (string, error) {
	lt := time.Duration(5 * time.Minute)
	exp := time.Now().Add(lt).Unix()
	createAt := time.Now().UTC()
	rawHash := member.Email + member.Username + fmt.Sprint(createAt) + "test"
	// authBody := AuthBody{
	// 	Data: domain.AuthMember{
	// 		UserID: member.Hash,
	// 		CreatedAt: createAt,
	// 	},
	// }
	authBody := AuthBody{
		Exp: exp,
		Data: domain.AuthMember{
			// Email:     member.Email,
			// Username:  member.Username,
			CreatedAt: createAt,
			HashAuth:  HashAuth(rawHash),
		},
	}

	bytesAuth, err := json.Marshal(authBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal member: %w", err)
	}

	encryptedMember, err := AESEncrypt(string(bytesAuth), "your-encryption-key-here")
	if err != nil {
		return "", fmt.Errorf("failed to encrypt member: %w", err)
	}

	token := encryptedMember

	return token, nil
}

func ValidateBearerToken(token string) (AuthBody, error) {
	authBody := AuthBody{}
	decrypted, err := AESDecrypt(token, "your-encryption-key-here")
	if err != nil {
		fmt.Println("Error decrypting token:", err)
		return authBody, fmt.Errorf("failed to decrypt token: %w", err)
	}

	if err = json.Unmarshal(decrypted, &authBody); err != nil {
		return authBody, fmt.Errorf("failed to unmarshal decrypted token: %w", err)
	}

	if authBody.Exp < time.Now().Unix() {
		return authBody, fmt.Errorf("token has expired")
	}

	return authBody, nil
}
