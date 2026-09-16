package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_HOST") + ":" + GetEnv("REDIS_PORT"),
		Password: GetEnv("REDIS_PASS"),
		DB:       0,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis not connected:", err)
	}

	log.Println("Redis connected")
}
