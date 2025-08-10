package server

import (
	"net/http"
	"time"

	ratelimiter "github.com/epicmet/rate-limit-algorithms/server/rate-limiter"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	r.GET(
		"/ping",
		ratelimiter.TokenBucket(
			ratelimiter.Config{
				BucketSize: 5,
				RefillRate: time.Second * 5,
				KeyPrefix: "ping",
			},
		),
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"message": "pong",
				},
			)
		})

	r.GET(
		"/idk",
		ratelimiter.TokenBucket(
			ratelimiter.Config{
				BucketSize: 5,
				RefillRate: time.Second * 10,
				KeyPrefix: "idk",
			},
		),
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"message": "wat?",
				},
			)
		})
	r.Run(":6969")
}
