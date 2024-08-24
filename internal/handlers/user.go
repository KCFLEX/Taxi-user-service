package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignUP(ctx context.Context, User models.UserInfo) error
}

type Handler struct {
	srv    Service
	router *gin.Engine
	port   string
}

func New(config config.Config, srv Service) *Handler {
	router := gin.Default()
	return &Handler{
		srv:    srv,
		router: router,
		port:   config.Port,
	}
}

func (h *Handler) RegisterRoutes() {

	h.router.Use(gin.Recovery())
	h.router.POST("/signup", h.SignUP)
}

func (h *Handler) Serve() error {

	h.RegisterRoutes()
	return h.router.Run(":" + h.port)
}

type SignINInfo struct {
}

func (h *Handler) SignUP(ctx *gin.Context) {

	var User models.UserInfo

	err := ctx.ShouldBind(&User)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.srv.SignUP(ctx, User)
	if err != nil {
		log.Printf("failed to signUP user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}
