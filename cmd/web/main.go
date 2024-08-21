package main

import (
	"log"

	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":" + config.Port)
}
