package main

import (
	"fmt"

	"app/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("test")

	r := gin.Default()
	r.GET("/", controllers.TopPage)
	r.Run()
}
