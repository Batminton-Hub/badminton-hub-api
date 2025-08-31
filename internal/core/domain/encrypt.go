package domain

import "time"

type EncrypBody struct {
	ByteBody []byte
	Exp      time.Time
}
