package service

import (
	"Badminton-Hub/internal/core/port"
	"time"
)

type CacheUtil struct {
	cache port.Cache
}

func NewCacheUtil(cache port.Cache) *CacheUtil {
	return &CacheUtil{
		cache: cache,
	}
}

func (c *CacheUtil) GetGoogleState(key string) (bool, error) {
	return c.cache.GetGoogleState(key)
}

func (c *CacheUtil) SetGoogleState(key string, lt time.Duration) error {
	return c.cache.SetGoogleState(key, lt)
}
