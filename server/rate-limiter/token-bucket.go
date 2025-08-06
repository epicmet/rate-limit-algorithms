package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/epicmet/rate-limit-algorithms/server/redis"
	"github.com/gin-gonic/gin"
)

type Config struct {
	BucketSize int
	RefillRate time.Duration
}

var ctx = context.Background()
var redisKey = "rate-limiter::token-bucket::bucket"

func TokenBucket(c Config) gin.HandlerFunc {

	rdb := redis.GetRedisClient()
	ticker := time.NewTicker(c.RefillRate)

	rdb.Set(ctx, redisKey, c.BucketSize, 0)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				{
					rdb.Set(ctx, redisKey, c.BucketSize, 0)
				}
			}
		}
	}()

	return func(c *gin.Context) {
		bucketCounterStr, err := redis.GetIntKey(ctx, redisKey)
		if err != nil || bucketCounterStr <= 0 {
			c.JSON(
				http.StatusTooManyRequests,
				gin.H{},
			)
			c.Abort()
		} else {
			if _, err := rdb.Decr(ctx, redisKey).Result(); err != nil {
				fmt.Println("Couldn't Decr the bucket counter. Error: ", err.Error())

				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"message": "Rate limiter faild",
					},
				)
				c.Abort()
			}
		}
	}
}
