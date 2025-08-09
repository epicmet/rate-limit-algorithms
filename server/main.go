package server

import (
	"net/http"
	"time"

	ratelimiter "github.com/epicmet/rate-limit-algorithms/server/rate-limiter"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	tokenBucketMiddleware :=
		ratelimiter.TokenBucket(
			ratelimiter.Config{
				BucketSize: 10,
				RefillRate: time.Second * 5,
			},
		)

	r := gin.Default()
	r.GET(
		"/ping",
		tokenBucketMiddleware,
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"message": "pong",
				},
			)
		})
	r.Run(":6969")
}
