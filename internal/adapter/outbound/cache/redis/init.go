package redisCache

import (
	"Badminton-Hub/util"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client redis.Client
}

func NewRedisCache() *RedisCache {
	config, err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisCacheAddr, // use default Addr
		Password: config.RedisCachePassword,
		DB:       config.RedisCacheDB,
		Protocol: 2,
	})
	return &RedisCache{
		client: *client,
	}
}
