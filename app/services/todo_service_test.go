package services

import (
	"app/dto"
	"app/models"
	"app/repositories"
	"app/test/factories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestTodoServiceSuite struct {
	WithDbSuite
}

var (
	user            *models.User
	testTodoService TodoService
)

func (s *TestTodoServiceSuite) SetupTest() {
	s.SetDbCon()

	// NOTE: テスト用ユーザの作成
	user = factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DbCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	todoRepository := repositories.NewTodoRepository(DbCon)
	testTodoService = NewTodoService(todoRepository)
}

func (s *TestTodoServiceSuite) TearDownTest() {
	s.CloseDb()
}

func (s *TestTodoServiceSuite) TestCreateTodo() {
	requestParams := dto.CreateTodoRequest{Title: "test title 1", Content: "test content 1"}

	result := testTodoService.CreateTodo(requestParams, user.ID)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)

	// NOTE: Todoリストが作成されていることを確認
	todo := models.Todo{}
	if err := DbCon.Where("user_id = ?", user.ID).First(&todo).Error; err != nil {
		s.T().Fatalf("failed to create todo %v", err)
	}
	assert.Equal(s.T(), "test title 1", todo.Title)
	assert.Equal(s.T(), "test content 1", todo.Content)
}

func (s *TestTodoServiceSuite) TestCreateTodo_ValidationError() {
	requestParams := dto.CreateTodoRequest{Title: "", Content: "test content 1"}

	result := testTodoService.CreateTodo(requestParams, user.ID)

	assert.NotNil(s.T(), result.Error)
	assert.Equal(s.T(), "validationError", result.ErrorType)

	// NOTE: Todoリストが作成されていないことを確認
	todo := models.Todo{}
	err := DbCon.Where("user_id = ?", user.ID).First(&todo).Error
	assert.NotNil(s.T(), err)
}

func (s *TestTodoServiceSuite) TestFetchTodosList() {
	testTodos := []models.Todo{
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
	if err := DbCon.Create(&testTodos).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	result := testTodoService.FetchTodosList(user.ID)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)
	assert.Len(s.T(), result.Todos, 2)
}

func (s *TestTodoServiceSuite) TestFetchTodo() {
	testTodo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DbCon.Create(&testTodo).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	result := testTodoService.FetchTodo(testTodo.ID, user.ID)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)
	assert.Equal(s.T(), testTodo.Title, result.Todo.Title)
}

func (s *TestTodoServiceSuite) TestUpdateTodo() {
	testTodo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DbCon.Create(&testTodo).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	requestParams := dto.UpdateTodoRequest{Title: "test updated title 1", Content: "test updated content 1"}
	result := testTodoService.UpdateTodo(testTodo.ID, requestParams, user.ID)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)
	assert.Equal(s.T(), "test updated title 1", result.Todo.Title)
	assert.Equal(s.T(), "test updated content 1", result.Todo.Content)
}

func (s *TestTodoServiceSuite) TestUpdateTodo_ValidationError() {
	testTodo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DbCon.Create(&testTodo).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	requestParams := dto.UpdateTodoRequest{Title: "", Content: "test updated content 1"}
	result := testTodoService.UpdateTodo(testTodo.ID, requestParams, user.ID)

	assert.NotNil(s.T(), result.Error)
	assert.Equal(s.T(), "validationError", result.ErrorType)
	// NOTE: Todoが更新されていないこと
	todo := models.Todo{}
	DbCon.Where("user_id = ?", user.ID).First(&todo)
	assert.Equal(s.T(), "test title 1", todo.Title)
	assert.Equal(s.T(), "test content 1", todo.Content)
}

func (s *TestTodoServiceSuite) TestDeleteTodo() {
	testTodo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DbCon.Create(&testTodo).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	result := testTodoService.DeleteTodo(testTodo.ID, user.ID)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)
}

func TestTodoService(t *testing.T) {
	// テストスイートを実行
	suite.Run(t, new(TestTodoServiceSuite))
}
