package main

import (
	"net/http"

	"github.com/blyndusk/cofy/api/database"
	"github.com/blyndusk/cofy/api/router"
	"github.com/gin-gonic/gin"
)

func main() {
	setupServer()
}

func setupServer() *gin.Engine {
	database.Connect()
	database.Migrate()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Setup(r)
	r.Run(":3003")
	return r
}
