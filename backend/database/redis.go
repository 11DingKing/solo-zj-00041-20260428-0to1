package database

import (
	"context"
	"fmt"
	"log"

	"canteen-system/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.GetAddr(),
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Redis connected successfully!")
}
