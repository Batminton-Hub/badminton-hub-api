package redisCache

import (
	"context"
	"time"
)

func (r *RedisCache) SetGoogleState(ctx context.Context, key string, lt time.Duration) error {
	if err := r.client.SetEx(ctx, key, true, lt).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) GetGoogleState(ctx context.Context, key string) (bool, error) {
	val, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

func (r *RedisCache) DelGoogleState(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
