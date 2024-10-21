package repositories

import (
	"app/models"
	"app/test/factories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var user *models.User

type TestTodoRePositorySuite struct {
	WithDbSuite
}

func (s *TestTodoRePositorySuite) SetupTest() {
	s.SetDbCon()

	// NOTE: テスト用ユーザの作成
	user = factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DbCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}
}

func (s *TestTodoRePositorySuite) TearDownTest() {
	s.CloseDb()
}

func (s *TestTodoRePositorySuite) TestCreateTodo() {
	insertTodo := models.Todo{}
	insertTodo.Title = "test title 1"
	insertTodo.Content = "test content 1"
	insertTodo.UserID = user.ID

	tr := NewTodoRepository(DbCon)
	err := tr.CreateTodo(&insertTodo)

	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), 0, insertTodo.ID)
}

func (s *TestTodoRePositorySuite) TestGetAllTodos() {
	insertTodos := []models.Todo{
		{
			Title:   "test title 1",
			Content: "test content 1",
			UserID:  user.ID,
		},
		{
			Title:   "test title 2",
			Content: "test content 2",
			UserID:  user.ID,
		},
	}
	if err := DbCon.Create(&insertTodos).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	todos := []models.Todo{}
	tr := NewTodoRepository(DbCon)
	tr.GetAllTodos(&todos, user.ID)

	assert.Equal(s.T(), 2, len(todos))
}

func (s *TestTodoRePositorySuite) TestGetTodoById() {
	insertTodo := models.Todo{}
	insertTodo.Title = "test title 1"
	insertTodo.Content = "test content 1"
	insertTodo.UserID = user.ID
	if err := DbCon.Create(&insertTodo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}

	todo := models.Todo{}
	tr := NewTodoRepository(DbCon)
	err := tr.GetTodoById(&todo, insertTodo.ID, user.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), insertTodo.ID, todo.ID)
}

func (s *TestTodoRePositorySuite) TestUpdateTodo() {
	todo := models.Todo{}
	todo.Title = "test title 1"
	todo.Content = "test content 1"
	todo.UserID = user.ID
	if err := DbCon.Create(&todo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}
	assert.Equal(s.T(), "test title 1", todo.Title)
	assert.Equal(s.T(), "test content 1", todo.Content)

	tr := NewTodoRepository(DbCon)
	todo.Title = "test updated title 1"
	todo.Content = "test updated content 1"
	err := tr.UpdateTodo(&todo)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test updated title 1", todo.Title)
	assert.Equal(s.T(), "test updated content 1", todo.Content)
}

func (s *TestTodoRePositorySuite) TestDeleteTodo() {
	todo := models.Todo{}
	todo.Title = "test title 1"
	todo.Content = "test content 1"
	todo.UserID = user.ID
	if err := DbCon.Create(&todo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}

	todos := []models.Todo{}
	tr := NewTodoRepository(DbCon)
	tr.GetAllTodos(&todos, user.ID)
	assert.Equal(s.T(), 1, len(todos))

	err := tr.DeleteTodo(&todo)

	assert.Nil(s.T(), err)
	tr.GetAllTodos(&todos, user.ID)
	assert.Equal(s.T(), 0, len(todos))
}

func TestTodoRepository(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestTodoRePositorySuite))
}
