package controllers

import (
	"app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService}
}

func (authController *authController) SignUp(ctx *gin.Context) {
	user, error := authController.authService.SignUp(ctx)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": error,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (authController *authController) SignIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "auth sign in controller",
	})
}
