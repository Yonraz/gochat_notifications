package initializers

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectToRedis() {
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr: host,
		Password: password,
		DB: 0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("could not connect to redis: %v\n", err)
	}

	fmt.Printf("Connected to redis!\n")
}