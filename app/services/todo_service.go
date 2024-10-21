package services

import (
	"app/dto"
	"app/models"
	"app/repositories"

	"github.com/go-playground/validator/v10"
)

type TodoService interface {
	CreateTodo(requestParams dto.CreateTodoRequest, userID int) *dto.CreateTodoResponse
	FetchTodosList(userID int) *dto.TodosListResponse
	FetchTodo(id int, userID int) *dto.FetchTodoResponse
	UpdateTodo(id int, requestParams dto.UpdateTodoRequest, userID int) *dto.UpdateTodoResponse
	DeleteTodo(id int, userID int) *dto.DeleteTodoResponse
}

type todoService struct {
	todoRepository repositories.TodoRepository
}

func NewTodoService(todoRepository repositories.TodoRepository) TodoService {
	return &todoService{todoRepository}
}

func (ts *todoService) CreateTodo(requestParams dto.CreateTodoRequest, userID int) *dto.CreateTodoResponse {
	todo := models.Todo{}
	todo.Title = requestParams.Title
	todo.Content = requestParams.Content
	todo.UserID = userID
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

func (ts *todoService) FetchTodosList(userID int) *dto.TodosListResponse {
	todos := []models.Todo{}
	error := ts.todoRepository.GetAllTodos(&todos, userID)
	if error != nil {
		return &dto.TodosListResponse{Todos: []models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	return &dto.TodosListResponse{Todos: todos, Error: nil, ErrorType: ""}
}

func (ts *todoService) FetchTodo(id int, userID int) *dto.FetchTodoResponse {
	todo := models.Todo{}
	error := ts.todoRepository.GetTodoByID(&todo, id, userID)
	if error != nil {
		return &dto.FetchTodoResponse{Todo: models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	return &dto.FetchTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) UpdateTodo(id int, requestParams dto.UpdateTodoRequest, userID int) *dto.UpdateTodoResponse {
	todo := models.Todo{}
	error := ts.todoRepository.GetTodoByID(&todo, id, userID)
	if error != nil {
		return &dto.UpdateTodoResponse{Todo: models.Todo{}, Error: error, ErrorType: "notFound"}
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
	if updateError != nil {
		return &dto.UpdateTodoResponse{Todo: todo, Error: updateError, ErrorType: "internalServerError"}
	}
	return &dto.UpdateTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) DeleteTodo(id int, userID int) *dto.DeleteTodoResponse {
	todo := models.Todo{}
	error := ts.todoRepository.GetTodoByID(&todo, id, userID)
	if error != nil {
		return &dto.DeleteTodoResponse{Error: error, ErrorType: "notFound"}
	}

	deleteError := ts.todoRepository.DeleteTodo(&todo)
	if deleteError != nil {
		return &dto.DeleteTodoResponse{Error: deleteError, ErrorType: "internalServerError"}
	}
	return &dto.DeleteTodoResponse{Error: nil, ErrorType: ""}
}
