package util

import (
	"Badminton-Hub/internal/core/domain"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

func GenerateHash(key string) string {
	timeNow := time.Now().UnixMilli()
	data := fmt.Sprint(key, timeNow)
	hashBy := sha256.New()
	hashBy.Write([]byte(data))
	bytesDigest := hashBy.Sum(nil)

	hash := fmt.Sprintf("%x", bytesDigest)
	return hash
}

func HashPassword(password, key string) string {
	data := fmt.Sprint(password, key)
	hashBy := sha256.New()
	hashBy.Write([]byte(data))
	bytesDigest := hashBy.Sum(nil)

	newPassword := fmt.Sprintf("%x", bytesDigest)
	return newPassword
}

func HashAuth(key string) string {
	data := fmt.Sprint(key + "hash_auth")
	hashBy := sha256.New()
	hashBy.Write([]byte(data))
	bytesDigest := hashBy.Sum(nil)

	newPassword := fmt.Sprintf("%x", bytesDigest)
	return newPassword
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

func AESEncrypt(plaintext string, key string) (string, error) {
	iv := "my16digitIvKey12"

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}
func AESDecrypt(encrypted, key string) ([]byte, error) {
	iv := "my16digitIvKey12"

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}

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
	authBody := AuthBody{
		Exp: exp,
		Data: domain.AuthMember{
			Email:     member.Email,
			Username:  member.Username,
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

	return authBody, nil
}
