package server

import (
	"net/http"
	"time"

	"github.com/epicmet/rate-limit-algorithms/server/rate-limiter/algorithms/token-bucket"
	statemanager "github.com/epicmet/rate-limit-algorithms/server/rate-limiter/state-manager"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()

	stateManager := statemanager.New(
		"redis",
		statemanager.Config{
			Addr: "localhost:6379",
		},
	)

	r.GET(
		"/ping",
		tokenbucket.New("ping", 5, time.Second*5, stateManager).GinMiddleware(),
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
		tokenbucket.New("idk", 5, time.Second*10, stateManager).GinMiddleware(),
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
