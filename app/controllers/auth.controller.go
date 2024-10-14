package controllers

import (
	"app/services"
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
	token, error := authController.authService.SignIn(ctx)
	if error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User Not Found",
		})
		return
	}

	// NOTE: Cookieにtokenをセット
	ctx.SetCookie("token", token, 3600*24, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
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
