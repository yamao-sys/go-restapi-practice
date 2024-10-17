package main

import (
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
	dbCon := db.Init()

	// repository
	userRepository := repositories.NewUserRepository(dbCon)
	todoRepository := repositories.NewTodoRepository(dbCon)

	// service
	authService := services.NewAuthService(userRepository)
	todoService := services.NewTodoService(todoRepository)

	// controller
	authController := controllers.NewAuthController(authService)
	todoController := controllers.NewTodoController(todoService, authService)
	authRouter := routers.NewAuthRouter(authController)
	todoRouter := routers.NewTodoRouter(todoController)

	// router
	r := gin.Default()

	r.GET("/", controllers.TopPage)
	authRouter.SetRouting(r)
	todoRouter.SetRouting(r)
	r.Run(":" + strconv.Itoa(config.Config.ServerPort))
}
