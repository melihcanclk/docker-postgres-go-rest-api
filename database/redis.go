package database

import (
	"context"
	"fmt"

	"github.com/melihcanclk/docker-postgres-go-rest-api/config"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         context.Context
)

func ConnectToRedis() {
	port := config.RedisPort

	ctx = context.Background()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	err := RedisClient.Set(ctx, "test", "How to Refresh Access Tokens the Right Way in Golang", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Redis client connected successfully...")

}
