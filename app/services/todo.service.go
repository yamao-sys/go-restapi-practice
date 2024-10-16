package services

import (
	"app/dto"
	"app/models"
	"app/repositories"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TodoService interface {
	CreateTodo(ctx *gin.Context, userId int) *dto.CreateTodoResponse
}

type todoService struct {
	todoRepository repositories.TodoRepository
}

func NewTodoService(todoRepository repositories.TodoRepository) TodoService {
	return &todoService{todoRepository}
}

func (ts *todoService) CreateTodo(ctx *gin.Context, userId int) *dto.CreateTodoResponse {
	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.CreateTodoRequest{}
	if err := ctx.ShouldBind(&requestParams); err != nil {
		return &dto.CreateTodoResponse{Todo: models.Todo{}, Error: err, ErrorType: "internalServerError"}
	}

	todo := models.Todo{}
	todo.Title = requestParams.Title
	todo.Content = requestParams.Content
	todo.UserID = userId
	// NOTE: バリデーションチェック
	validate := validator.New()
	validationErrors := validate.Struct(todo)
	if validationErrors != nil {
		return &dto.CreateTodoResponse{Todo: todo, Error: validationErrors, ErrorType: "validationError"}
	}

	// NOTE: Create処理
	err := ts.todoRepository.CreateTodo(&todo)
	if err != nil {
		return &dto.CreateTodoResponse{Todo: todo, Error: err, ErrorType: "internalServerError"}
	}
	return &dto.CreateTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}
