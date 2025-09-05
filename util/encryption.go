package util

import (
	"Badminton-Hub/internal/core/domain"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// One way Encryption
func Sha256(data any) string {
	hashBy := sha256.New()
	hashBy.Write([]byte(fmt.Sprint(data)))
	bytesDigest := hashBy.Sum(nil)

	hash := fmt.Sprintf("%x", bytesDigest)
	return hash
}

func MD5(data any) string {
	hashBy := md5.New()
	hashBy.Write([]byte(fmt.Sprint(data)))
	bytesDigest := hashBy.Sum(nil)

	hash := fmt.Sprintf("%x", bytesDigest)
	return hash
}

// Two way Encryption
func AESEncrypt(body any, key string, lt time.Duration) (string, error) {
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

	iv, err := randomIV()
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainTextBlock)

	ciphertext = append(iv, ciphertext...) // แนบ iv ไปกับ ciphertext
	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}
func AESDecrypt(encryptData string, key string, body any) error {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return err
	}

	// แยก iv ออกจาก ciphertext
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
	if err := DecryptGOB(ciphertext, &payload); err != nil {
		return err
	}

	if err := DecryptGOB(payload.ByteBody, body); err != nil {
		return err
	}

	return nil
}

func JWTEncrypt(body any, key string, lt time.Duration) (string, error) {
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
func JWTDecrypt(encryptData string, key string, body any) error {
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
		if err := DecryptGOB(byteBody, body); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// Byte Convert
func EncryptGOB(body any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func DecryptGOB(data []byte, body any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(body); err != nil {
		return err
	}
	return nil
}

// other function
func randomIV() ([]byte, error) {
	config := LoadConfig()

	if config.Mode == "DEVERLOP" {
		return config.DefaultAESIV, nil
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return []byte{}, err
	}
	return iv, nil
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
