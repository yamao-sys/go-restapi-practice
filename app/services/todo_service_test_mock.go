package services

import (
	"app/dto"
	"app/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TodoServiceTestSuite struct {
	suite.Suite
}

type MockTodoRepository struct {
	mock.Mock
}

func (_m *MockTodoRepository) CreateTodo(todo *models.Todo) error {
	ret := _m.Called(todo)
	return ret.Error(0)
}

func (_m *MockTodoRepository) GetAllTodos(todos *[]models.Todo, userID int) error {
	ret := _m.Called(todos, userID)
	return ret.Error(0)
}

func (_m *MockTodoRepository) GetTodoByID(todo *models.Todo, id int, userID int) error {
	ret := _m.Called(todo, id, userID)
	return ret.Error(0)
}

func (_m *MockTodoRepository) UpdateTodo(todo *models.Todo) error {
	ret := _m.Called(todo)
	return ret.Error(0)
}

func (_m *MockTodoRepository) DeleteTodo(todo *models.Todo) error {
	ret := _m.Called(todo)
	return ret.Error(0)
}

func (s *TodoServiceTestSuite) TestCreateTodo() {
	// todoRepositoryをmock化
	mockTodoRepository := new(MockTodoRepository)
	mockTodoRepository.On("CreateTodo", &models.Todo{Title: "test title 1", Content: "test content 1", UserID: 1}).Return(nil)

	ts := NewTodoService(mockTodoRepository)
	result := ts.CreateTodo(dto.CreateTodoRequest{Title: "test title 1", Content: "test content 1"}, 1)

	assert.Equal(s.T(), nil, result.Error)
	assert.Equal(s.T(), "", result.ErrorType)
	assert.Equal(s.T(), "test title 1", result.Todo.Title)
	assert.Equal(s.T(), "test content 1", result.Todo.Content)
}

func (s *TodoServiceTestSuite) TestFetchTodosList() {
	// todoRepositoryをmock化
	mockTodoRepository := new(MockTodoRepository)
	mockTodoRepository.On("GetAllTodos", &[]models.Todo{}, 1).Return(nil)

	ts := NewTodoService(mockTodoRepository)
	result := ts.FetchTodosList(1)

	assert.Equal(s.T(), nil, result.Error)
	assert.Equal(s.T(), "", result.ErrorType)
}

func TestTodoServiceMock(t *testing.T) {
	suite.Run(t, new(TodoServiceTestSuite))
}
