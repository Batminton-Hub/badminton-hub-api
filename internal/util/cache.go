package core_util

import (
	"Badminton-Hub/internal/core/port"
	"context"
	"time"
)

type CacheUtil struct {
	cache port.CacheUtil
}

func NewCacheUtil(cache port.CacheUtil) *CacheUtil {
	return &CacheUtil{
		cache: cache,
	}
}

func (c *CacheUtil) GetGoogleState(ctx context.Context, key string) (bool, error) {
	return c.cache.GetGoogleState(ctx, key)
}

func (c *CacheUtil) SetGoogleState(ctx context.Context, key string, lt time.Duration) error {
	return c.cache.SetGoogleState(ctx, key, lt)
}

func (c *CacheUtil) DelGoogleState(ctx context.Context, key string) error {
	return c.cache.DelGoogleState(ctx, key)
}
