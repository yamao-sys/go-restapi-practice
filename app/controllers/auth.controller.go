package controllers

import (
	"app/services"
	"app/utils"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	GetUser(ctx *gin.Context)
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

func (authController *authController) GetUser(ctx *gin.Context) {
	// NOTE: Cookieからtokenを取得
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": "fail get cookie",
		})
		return
	}
	// NOTE: tokenに該当するユーザを取得する
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("abcdefghijklmn"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// NOTE: user_idに該当するユーザを取得して返す
		userId := int(claims["user_id"].(float64))
		user := authController.authService.Getuser(userId)
		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	} else {
		fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"error": "fail get claim",
		})
	}
}
