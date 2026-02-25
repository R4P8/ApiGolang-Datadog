package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

var RedisClient *redis.Client

func InitRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	redistrace.WrapClient(rdb,
		redistrace.WithServiceName("task-manager-redis"),
	)

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	RedisClient = rdb
	log.Println(" Connected to Redis with Datadog tracing")

	return rdb, nil
}
