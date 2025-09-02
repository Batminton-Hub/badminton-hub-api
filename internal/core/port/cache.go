package port

import "time"

type Cache interface {
	GetGoogleState(key string) (bool, error)
	SetGoogleState(key string, lt time.Duration) error
}
