package redisCache

import (
	"context"
	"time"
)

func (r *RedisCache) SetGoogleState(key string, lt time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.SetEx(ctx, key, true, lt).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) GetGoogleState(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

func (r *RedisCache) DelGoogleState(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
