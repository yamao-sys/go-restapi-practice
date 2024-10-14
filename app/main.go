package main

import (
	"fmt"
	"strconv"

	"app/config"
	"app/controllers"
	"app/db"
	"app/repositories"
	"app/routers"
	"app/services"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("test")

	dbCon := db.Init()

	// repository
	userRepository := repositories.NewUserRepository(dbCon)

	// service
	authService := services.NewAuthService(userRepository)

	// controller
	authController := controllers.NewAuthController(authService)
	authRouter := routers.NewAuthRouter(authController)

	// router
	r := gin.Default()

	r.GET("/", controllers.TopPage)
	authRouter.SetRouting(r)
	r.Run(":" + strconv.Itoa(config.Config.ServerPort))
}
