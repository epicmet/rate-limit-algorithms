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
				BucketSize: 10,
				RefillRate: time.Second * 5,
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
	r.Run(":6969")
}
