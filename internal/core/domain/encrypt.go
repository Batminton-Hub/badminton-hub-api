package domain

import "time"

type AESEncryption struct{}

type JWTEncryption struct{}

type EncrypBody struct {
	ByteBody []byte
	Exp      time.Time
}
