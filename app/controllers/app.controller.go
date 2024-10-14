package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TopPage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello Gin!!!!",
	})
}
