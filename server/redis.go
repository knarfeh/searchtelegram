package main

import (
	redis "github.com/go-redis/redis"
	"github.com/knarfeh/searchtelegram/server/diagnose"
	"time"
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

// Diagnose start diagnose check
func (rClient *RedisClient) Diagnose() diagnose.ComponentReport {
	var (
		err   error
		start time.Time
	)

	report := diagnose.NewReport("redis")
	start = time.Now()
	err = rClient.Client.Ping().Err()
	report.AddLatency(start)
	report.Check(err, "Redis client ping failed", "Check environment variables or redis health")
	start = time.Now()
	return *report
}
