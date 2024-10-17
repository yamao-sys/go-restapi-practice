package dto

import "app/models"

type CreateTodoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateTodoResponse struct {
	Todo      models.Todo
	Error     error
	ErrorType string
}

type UpdateTodoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateTodoResponse struct {
	Todo      models.Todo
	Error     error
	ErrorType string
}
