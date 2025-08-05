package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "pong",
			},
		)
	})
	r.Run(":6969")
}
