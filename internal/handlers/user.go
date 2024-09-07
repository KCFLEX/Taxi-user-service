package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignUP(ctx context.Context, User models.UserInfo) error
	SignIN(ctx context.Context, user models.UserInfo) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	CheckTokenInRedis(ctx context.Context, token string) error
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
	h.router.POST("/signin", h.SignIN)

	h.router.Use(h.AuthMiddleWare)

	// Now all routes defined after this line will be protected
	h.router.POST("/logout", h.LogOut)
	h.router.GET("/profile", h.Profile) //
}

func (h *Handler) Serve() error {

	h.RegisterRoutes()
	return h.router.Run(":" + h.port)
}

type SignINInfo struct {
	Phone    string
	Password string
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

func (h *Handler) SignIN(ctx *gin.Context) {

	var userCred SignINInfo

	err := ctx.ShouldBind(&userCred)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	user := models.UserInfo{
		PhoneNO:  userCred.Phone,
		Password: userCred.Password,
	}

	tokenStr, err := h.srv.SignIN(ctx, user)
	if err != nil {
		if errors.Is(err, errorpac.ErrUserDoesNotExist) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User does not exist"})
		} else if errors.Is(err, errorpac.ErrPasswordInvalid) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	//fmt.Println("this is the token", tokenStr)
	// stores token in cookies
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(ctx.Writer, &cookie)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "user successfully authorized"})
}

func (h *Handler) LogOut(ctx *gin.Context) {
	cookie := &http.Cookie{
		Name:     "auth_token",         // The name of the cookie to clear
		Value:    "",                   // Empty value for clearing
		Path:     "/",                  // Path must match the original cookie's Path
		MaxAge:   -1,                   // Negative MaxAge deletes the cookie
		HttpOnly: true,                 // Same HttpOnly setting as the original cookie
		Secure:   true,                 // Same Secure setting as the original cookie
		SameSite: http.SameSiteLaxMode, // Same SameSite setting as the original cookie
	}
	fmt.Println("------logout")
	http.SetCookie(ctx.Writer, cookie)

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *Handler) AuthMiddleWare(ctx *gin.Context) {
	// Retrieve the token from the cookie
	tokenStr, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
		return
	}
	fmt.Println("-------")
	fmt.Print(tokenStr)

	//verify token
	_, err = h.srv.VerifyToken(ctx, tokenStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		ctx.Abort()
		return
	}

	fmt.Println("-------")
	err = h.srv.CheckTokenInRedis(ctx, tokenStr)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		ctx.Abort()
	}
	fmt.Println("--------")
	ctx.Next()

}

func (h *Handler) Profile(ctx *gin.Context) {
	// This is a protected route. Only accessible if the JWT token is valid.
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected profile area!"})
}
