package util

import (
	"crypto/sha256"
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
	timeNow := time.Now().UnixMilli()
	data := fmt.Sprint(password, key, timeNow)
	hashBy := sha256.New()
	hashBy.Write([]byte(data))
	bytesDigest := hashBy.Sum(nil)

	newPassword := fmt.Sprintf("%x", bytesDigest)
	return newPassword
}

func AESEncrypt(plaintext, key string) (string, error) {
	// Placeholder for AES encryption logic
	return "", nil
}
func AESDecrypt(ciphertext, key string) (string, error) {
	// Placeholder for AES decryption logic
	return "", nil
}
