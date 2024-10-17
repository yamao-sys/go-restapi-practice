package controllers

import (
	"app/services"
	"app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodoController interface {
	Create(ctx *gin.Context)
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type todoController struct {
	todoService services.TodoService
	authService services.AuthService
}

func NewTodoController(todoService services.TodoService, authService services.AuthService) TodoController {
	return &todoController{todoService, authService}
}

func (todoController *todoController) Create(ctx *gin.Context) {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized error",
		})
		return
	}

	result := todoController.todoService.CreateTodo(ctx, user.ID)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"todo": result.Todo,
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

func (todoController *todoController) Index(ctx *gin.Context) {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized error",
		})
		return
	}

	result := todoController.todoService.FetchTodosList(ctx, user.ID)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"todos": result.Todos,
		})
		return
	}

	switch result.ErrorType {
	case "notFound":
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
	case "internalServerError":
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
	}
}

func (todoController *todoController) Show(ctx *gin.Context) {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized error",
		})
		return
	}

	result := todoController.todoService.FetchTodo(ctx, user.ID)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"todo": result.Todo,
		})
		return
	}

	switch result.ErrorType {
	case "notFound":
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
	case "internalServerError":
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
	}
}

func (todoController *todoController) Update(ctx *gin.Context) {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized error",
		})
		return
	}

	result := todoController.todoService.UpdateTodo(ctx, user.ID)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"todo": result.Todo,
		})
		return
	}

	switch result.ErrorType {
	case "validationError":
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": utils.CoordinateValidationErrors(result.Error),
		})
	case "notFound":
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
	case "internalServerError":
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
	}
}

func (todoController *todoController) Delete(ctx *gin.Context) {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized error",
		})
		return
	}
	result := todoController.todoService.DeleteTodo(ctx, user.ID)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": "delete todo(ID: " + ctx.Param("id") + ") successfully",
		})
		return
	}

	switch result.ErrorType {
	case "notFound":
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
	case "internalServerError":
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
	}
}
