package port

import "time"

type Cache interface {
	CacheGoogle
}

type CacheGoogle interface {
	GetGoogleState(key string) (bool, error)
	SetGoogleState(key string, lt time.Duration) error
	DelGoogleState(key string) error
}
