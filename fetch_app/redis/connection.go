package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func Initiate() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("error initializing Redis Client: %v\n", err)
	}
	fmt.Println("⇨ Redis Client Started !!")
}
