package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"fmt"
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
	// rawHash := authBody.Data.UserID + fmt.Sprint(authBody.Data.CreatedAt) + key
	data := fmt.Sprint(rawHash + "hash_auth")
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

func EncrypteJwt(body any, key string, lt time.Duration) (string, error) {
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

func DecrypteJwt(encryptData string, key string, body any) error {
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
	if byteBody, err := encryptGOB(rawBody); err == nil {
		if err := decrypteGOB(byteBody, body); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func encryptGOB(body any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decrypteGOB(data []byte, body any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(body); err != nil {
		return err
	}
	return nil
}
