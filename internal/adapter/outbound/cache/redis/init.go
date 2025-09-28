package redis

import (
	"Badminton-Hub/util"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type CloseRedisCache func()

type RedisCache struct {
	client redis.Client
}

func NewRedisCache() (*RedisCache, CloseRedisCache) {
	config := util.LoadConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisCacheAddr, // use default Addr
		Password: config.RedisCachePassword,
		DB:       config.RedisCacheDB,
		Protocol: 2, // use default Protocol
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	redisCache := &RedisCache{
		client: *client,
	}

	closeRedis := closeRedisCache(client)

	return redisCache, closeRedis
}

func closeRedisCache(client *redis.Client) CloseRedisCache {
	return func() {
		fmt.Println("closeRedisCache client")
		client.Close()
	}
}
