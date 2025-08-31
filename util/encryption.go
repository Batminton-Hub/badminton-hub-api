package util

import (
	"Badminton-Hub/internal/core/domain"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func HashAuth(rawHash string) string {
	data := fmt.Sprint(rawHash + "hash_auth")
	hashBy := sha256.New()
	hashBy.Write([]byte(data))
	bytesDigest := hashBy.Sum(nil)

	newPassword := fmt.Sprintf("%x", bytesDigest)
	return newPassword
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

func aesEncrypt(body any, key string, lt time.Duration) (string, error) {
	bytePayload, err := EncryptGOB(body)
	if err != nil {
		return "", errors.New("failed to encrypt payload")
	}

	encryptBody := domain.EncrypBody{
		ByteBody: bytePayload,
		Exp:      time.Now().Add(lt),
	}

	payload, err := EncryptGOB(encryptBody)
	if err != nil {
		return "", errors.New("failed to encrypt payload")
	}

	var plainTextBlock []byte
	plaintext := string(payload)
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
		return "", errors.New("failed to create AES cipher")
	}

	iv := randomIV()
	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainTextBlock)

	ciphertext = append(iv, ciphertext...)
	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func aesDecrypt(encryptData string, key string, body any) error {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = pkcs5UnPadding(ciphertext)

	payload := domain.EncrypBody{}
	if err := DecrypteGOB(ciphertext, &payload); err != nil {
		return err
	}

	if err := DecrypteGOB(payload.ByteBody, body); err != nil {
		return err
	}

	return nil
}

func randomIV() []byte {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}
	return iv
}

func jwtEncrypt(body any, key string, lt time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"body": body,
		"exp":  time.Now().Add(lt).Unix(),
	}
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encryptData, err := tokenJWT.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return encryptData, nil
}

func jwtDecrypt(encryptData string, key string, body any) error {
	var tokenJWT *jwt.Token
	var err error
	tokenJWT, err = jwt.Parse(encryptData, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return err
	}

	exp, err := tokenJWT.Claims.GetExpirationTime()
	switch {
	case err != nil:
		return err
	case time.Now().After(exp.Time):
		return fmt.Errorf("token has expired")
	}

	rawBody := tokenJWT.Claims.(jwt.MapClaims)["body"]
	if byteBody, err := EncryptGOB(rawBody); err == nil {
		if err := DecrypteGOB(byteBody, body); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func EncryptGOB(body any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecrypteGOB(data []byte, body any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(body); err != nil {
		return err
	}
	return nil
}

type Encryptetion interface {
	Encrypte(body any, key string, lt time.Duration) (string, error)
	Decrypte(encryptData string, key string, body any) error
}

type JWTEncryption struct{}
type AESEncryption struct{}

func NewJWTEncryption() JWTEncryption { return JWTEncryption{} }
func NewAESEncryption() AESEncryption { return AESEncryption{} }

func (jwt *JWTEncryption) Encrypte(body any, key string, lt time.Duration) (string, error) {
	return jwtEncrypt(body, key, lt)
}
func (jwt *JWTEncryption) Decrypte(encryptData string, key string, body any) error {
	return jwtDecrypt(encryptData, key, body)
}

func (aes *AESEncryption) Encrypte(body any, key string, lt time.Duration) (string, error) {
	return aesEncrypt(body, key, lt)
}
func (aes *AESEncryption) Decrypte(encryptData string, key string, body any) error {
	return aesDecrypt(encryptData, key, body)
}
