package dto

import "app/models"

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	User      models.User
	Error     error
	ErrorType string
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	TokenString     string
	NotFoundMessage string
	Error           error
}
