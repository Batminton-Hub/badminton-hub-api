package service

import (
	"Badminton-Hub/util"
	"time"
)

type JWTEncryption struct{}
type AESEncryption struct{}

func NewJWTEncryption() *JWTEncryption { return &JWTEncryption{} }
func NewAESEncryption() *AESEncryption { return &AESEncryption{} }

func (jwt *JWTEncryption) Encrypte(body any, key string, lt time.Duration) (string, error) {
	return util.JWTEncrypt(body, key, lt)
}
func (jwt *JWTEncryption) Decrypte(encryptData string, key string, body any) error {
	return util.JWTDecrypt(encryptData, key, body)
}

func (aes *AESEncryption) Encrypte(body any, key string, lt time.Duration) (string, error) {
	return util.AESEncrypt(body, key, lt)
}
func (aes *AESEncryption) Decrypte(encryptData string, key string, body any) error {
	return util.AESDecrypt(encryptData, key, body)
}
