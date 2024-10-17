package routers

import (
	"app/controllers"

	"github.com/gin-gonic/gin"
)

type TodoRouter interface {
	SetRouting(r *gin.Engine)
}

type todoRouter struct {
	todoController controllers.TodoController
}

func NewTodoRouter(todoController controllers.TodoController) TodoRouter {
	return &todoRouter{todoController}
}

func (tr *todoRouter) SetRouting(r *gin.Engine) {
	r.POST("/todos/", tr.todoController.CreateTodo)
	r.GET("/todos/:id", tr.todoController.ShowTodo)
	r.PUT("/todos/:id", tr.todoController.UpdateTodo)
}
