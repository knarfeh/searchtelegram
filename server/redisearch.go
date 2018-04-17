package main

import (
	"fmt"

	"github.com/RedisLabs/redisearch-go/redisearch"
	"github.com/knarfeh/searchtelegram/server/diagnose"
	"time"
)

// RedisearchClient ...
type RedisearchClient struct {
	Client *redisearch.Client
}

func GetRedisearchSchema() *redisearch.Schema {
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("desc")).
		AddField(redisearch.NewTextField("type")).
		AddField(redisearch.NewTextField("title")).
		AddField(redisearch.NewTextField("tgid")).
		AddField(redisearch.NewTextField("tagsforsearch")).
		AddField(redisearch.NewTagFieldOptions("tags", redisearch.TagFieldOptions{Separator: '#'}))
	return sc
}

// NewRedisearchClient ...
func NewRedisearchClient(host string, port string) *RedisearchClient {
	client := redisearch.NewClient(host+":"+port, "st_index")
	// client.Drop()

	sc := GetRedisearchSchema()

	// Create the index with the given schema
	if err := client.CreateIndex(sc); err != nil {
		fmt.Print(err)
	}

	return &RedisearchClient{
		Client: client,
	}
}

// Diagnose start diagnose check
func (redisearchClient *RedisearchClient) Diagnose() diagnose.ComponentReport {
	var (
		start time.Time
	)

	report := diagnose.NewReport("redisearch")
	start = time.Now()
	_, err := redisearchClient.Client.Info()
	report.AddLatency(start)
	report.Check(err, "Redis client ping failed", "Check environment variables or redis health")
	start = time.Now()
	return *report
}
