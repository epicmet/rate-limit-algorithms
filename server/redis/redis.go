package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func GetRedisClient() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Couldn't connecto to Redis %v", err)
	}

	return rdb
}

func GetIntKey(ctx context.Context, key string) (int, error) {
	keyStr, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	keyInt, err := strconv.Atoi(keyStr)
	if err != nil {
		return 0, err
	}

	return keyInt, nil
}
