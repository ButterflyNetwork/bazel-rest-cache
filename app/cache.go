package app

import (
	"github.com/go-redis/redis"
	"time"
)

type BazelCache interface {
	Get(key string) (string, bool)
	Set(key string, value string) bool
}

type redisBazelCache struct {
	client *redis.Client
}

func NewRedisBazelCache(redisAddr string) BazelCache {
	return &redisBazelCache{
		client: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
			IdleTimeout: time.Second * 5,
			IdleCheckFrequency: time.Second * 5,
		}),
	}
}

func (c *redisBazelCache) Get(key string) (string, bool) {
	bytes, err := c.client.Get(key).Result()
	
	if err != nil {
		return "", false
	}
	
	return bytes, true
}

func (c *redisBazelCache) Set(key string, value string) bool {
	err := c.client.Set(key, value, 0).Err()
	return err == nil
}
