package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AddNewWallet(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
		return
	}

	//verify token
	userIDstr, err := h.srv.VerifyToken(ctx, tokenStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		ctx.Abort()
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fmt.Println(userID)
	var walletInfo models.Wallet

	err = ctx.ShouldBind(&walletInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("---------")
	err = h.srv.AddPersonalWallet(ctx, userID, walletInfo)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user personal wallet succesfully created"})
}

func (h *Handler) AddfamilyWallet(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
		return
	}

	//verify token
	userIDstr, err := h.srv.VerifyToken(ctx, tokenStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		ctx.Abort()
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fmt.Println(userID)

	var walletInfo models.Wallet

	err = ctx.ShouldBind(&walletInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.srv.AddFamilyWallet(ctx, userID, walletInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user family wallet successfully created"})

}

func (h *Handler) AddUserToFamilyWallet(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
		return
	}

	//verify token
	userIDstr, err := h.srv.VerifyToken(ctx, tokenStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		ctx.Abort()
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(userID)

	var phone models.Phone

	err = ctx.ShouldBind(&phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.srv.AddUserToFamilyByPhone(ctx, userID, phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "new member has been succesfully added to the family wallet"})

}
