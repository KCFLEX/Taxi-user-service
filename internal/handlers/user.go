package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignUP(ctx context.Context, User *models.SignUPInfo) error
}

type Handler struct {
	router *gin.Engine
}

func New() *Handler {
	router := gin.Default()
	return &Handler{
		router: router,
	}
}

func (h *Handler) SignUP(ctx *gin.Context) {

	var User models.SignUPInfo

	err := ctx.ShouldBind(&User)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// form validation
	errors := User.Required()

	if len(errors) > 0 {
		for field, message := range errors {
			fmt.Fprintf(ctx.Writer, "error in %s: %s\n", field, message)
			return
		}
	}

	if !User.Validate() {
		fmt.Fprint(ctx.Writer, User.Errors)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, User)

}
