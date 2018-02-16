package main

import (
	redis "github.com/go-redis/redis"
)

// RedisClient ...
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient ...
func NewRedisClient(host string, port string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})

	return &RedisClient{
		Client: client,
	}
}
