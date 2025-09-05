package port

import (
	"context"
	"time"
)

type Cache interface {
	CacheGoogle
}

type CacheGoogle interface {
	GetGoogleState(ctx context.Context, key string) (bool, error)
	SetGoogleState(ctx context.Context, key string, lt time.Duration) error
	DelGoogleState(ctx context.Context, key string) error
}
