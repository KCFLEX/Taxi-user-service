package handlers

import (
	"net/http"

	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) OrderTaxi(ctx *gin.Context) {
	var Order models.OrderInfo

	err := ctx.ShouldBind(&Order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	tokenStr, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
		return
	}

	userID, err := h.srv.VerifyToken(ctx, tokenStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		ctx.Abort()
		return
	}

	Order.UserID = userID

}
