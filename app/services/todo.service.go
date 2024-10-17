package services

import (
	"app/dto"
	"app/models"
	"app/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TodoService interface {
	CreateTodo(ctx *gin.Context, userId int) *dto.CreateTodoResponse
	FetchTodosList(ctx *gin.Context, userId int) *dto.TodosListResponse
	FetchTodo(ctx *gin.Context, userId int) *dto.FetchTodoResponse
	UpdateTodo(ctx *gin.Context, userId int) *dto.UpdateTodoResponse
	DeleteTodo(ctx *gin.Context, userId int) *dto.DeleteTodoResponse
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

func (ts *todoService) FetchTodosList(ctx *gin.Context, userId int) *dto.TodosListResponse {
	todos := []models.Todo{}
	error := ts.todoRepository.GetAllTodos(&todos, userId)
	if error != nil {
		return &dto.TodosListResponse{Todos: []models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	return &dto.TodosListResponse{Todos: todos, Error: nil, ErrorType: ""}
}

func (ts *todoService) FetchTodo(ctx *gin.Context, userId int) *dto.FetchTodoResponse {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return &dto.FetchTodoResponse{Todo: models.Todo{}, Error: err, ErrorType: "internalServerError"}
	}

	todo := models.Todo{}
	error := ts.todoRepository.GetTodoById(&todo, id, userId)
	if error != nil {
		return &dto.FetchTodoResponse{Todo: models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	return &dto.FetchTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) UpdateTodo(ctx *gin.Context, userId int) *dto.UpdateTodoResponse {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return &dto.UpdateTodoResponse{Todo: models.Todo{}, Error: err, ErrorType: "internalServerError"}
	}

	todo := models.Todo{}
	error := ts.todoRepository.GetTodoById(&todo, id, userId)
	if error != nil {
		return &dto.UpdateTodoResponse{Todo: models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.UpdateTodoRequest{}
	if err := ctx.ShouldBind(&requestParams); err != nil {
		return &dto.UpdateTodoResponse{Todo: models.Todo{}, Error: err, ErrorType: "internalServerError"}
	}
	todo.Title = requestParams.Title
	todo.Content = requestParams.Content
	// NOTE: バリデーションチェック
	validate := validator.New()
	validationErrors := validate.Struct(todo)
	if validationErrors != nil {
		return &dto.UpdateTodoResponse{Todo: todo, Error: validationErrors, ErrorType: "validationError"}
	}

	// NOTE: Update処理
	updateError := ts.todoRepository.UpdateTodo(&todo)
	if err != nil {
		return &dto.UpdateTodoResponse{Todo: todo, Error: updateError, ErrorType: "internalServerError"}
	}
	return &dto.UpdateTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) DeleteTodo(ctx *gin.Context, userId int) *dto.DeleteTodoResponse {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return &dto.DeleteTodoResponse{Error: err, ErrorType: "internalServerError"}
	}

	todo := models.Todo{}
	error := ts.todoRepository.GetTodoById(&todo, id, userId)
	if error != nil {
		return &dto.DeleteTodoResponse{Error: error, ErrorType: "notFound"}
	}

	deleteError := ts.todoRepository.DeleteTodo(&todo)
	if deleteError != nil {
		return &dto.DeleteTodoResponse{Error: err, ErrorType: "internalServerError"}
	}
	return &dto.DeleteTodoResponse{Error: nil, ErrorType: ""}
}
