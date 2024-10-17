package controllers

import (
	"app/services"
	"app/utils"
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
	result := authController.authService.SignUp(ctx)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"user": result.User,
		})
		return
	}

	switch result.ErrorType {
	case "internalServerError":
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
	case "validationError":
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": utils.CoordinateValidationErrors(result.Error),
		})
	}
}

func (authController *authController) SignIn(ctx *gin.Context) {
	result := authController.authService.SignIn(ctx)

	if result.NotFoundMessage != "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": result.NotFoundMessage,
		})
		return
	}
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	// NOTE: Cookieにtokenをセット
	ctx.SetCookie("token", result.TokenString, 3600*24, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"token": result.TokenString,
	})
}
