package handlers

import (
	"log"

	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes() {

	h.router.Use(gin.Recovery())
	h.router.POST("/signup", h.SignUP)
}

func (h *Handler) Serve() error {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	h.RegisterRoutes()
	return h.router.Run(":" + config.Port)
}
