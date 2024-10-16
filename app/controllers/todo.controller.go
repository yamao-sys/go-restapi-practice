package controllers

import (
	"github.com/gin-gonic/gin"
)

type TodoController interface {
	//
}

type todoController struct {
	//
}

func NewTodoController() TodoController {
	return &todoController{}
}

func (todoController *todoController) GetAllTodos(ctx *gin.Context) {
	//
}
