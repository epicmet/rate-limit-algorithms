package statemanager

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStruct struct {
	client *redis.Client
}

var rdb *redis.Client

func newRedis(addr string) *redisStruct {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   0,
		})
	}

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Couldn't connecto to Redis %v", err)
	}

	return &redisStruct{
		client: rdb,
	}
}

func (r redisStruct) Set(key string, value interface{}, expireTime time.Duration) (string, error) {
	return r.client.Set(context.Background(), key, value, expireTime).Result()
}

func (r redisStruct) GetIntValue(key string) (int64, error) {
	keyStr, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	keyInt, err := strconv.ParseInt(keyStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return keyInt, nil
}

func (r redisStruct) Decr(key string) (int64, error) {
	return r.client.Decr(context.Background(), key).Result()
}

func (r redisStruct) Incr(key string) (int64, error) {
	return r.client.Incr(context.Background(), key).Result()
}
