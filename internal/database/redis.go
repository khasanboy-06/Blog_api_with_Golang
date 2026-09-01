package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis(){
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr ==  ""{
		redisAddr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "",
		DB: 0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Redis ga ulanishda xatolik: %v", err)
	}

	fmt.Println("Redisga muvaffaqiyatli ulandi")
}