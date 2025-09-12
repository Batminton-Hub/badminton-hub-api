package port

import "time"

type EncryptionUtil interface {
	Encrypte(body any, key string, lt time.Duration) (string, error)
	Decrypte(encryptData string, key string, body any) error
}
