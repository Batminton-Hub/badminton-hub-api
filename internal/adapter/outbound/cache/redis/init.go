package redis

import (
	"Badminton-Hub/util"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client redis.Client
}

func NewRedisCache() *RedisCache {
	config := util.LoadConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisCacheAddr, // use default Addr
		Password: config.RedisCachePassword,
		DB:       config.RedisCacheDB,
		Protocol: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisCache{
		client: *client,
	}
}
