package controllers

import (
	"app/models"
	"app/repositories"
	"app/services"
	"app/test/factories"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	user               *models.User
	testTodoController TodoController
)

type TestTodoControllerSuite struct {
	WithDBSuite
}

func (s *TestTodoControllerSuite) SetupTest() {
	s.SetDBCon()

	// NOTE: テスト用ユーザの作成
	user = factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DBCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	userRepository := repositories.NewUserRepository(DBCon)
	todoRepository := repositories.NewTodoRepository(DBCon)

	authService := services.NewAuthService(userRepository)
	todoService := services.NewTodoService(todoRepository)

	// NOTE: テスト対象のコントローラを設定
	testTodoController = NewTodoController(todoService, authService)

	// NOTE: ログインし、tokenに値を格納
	s.SignIn()
}

func (s *TestTodoControllerSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestTodoControllerSuite) TestCreateTodo() {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	createTodoBody := bytes.NewBufferString("{\"title\":\"test title 1\",\"content\":\"test content 1\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/todos", createTodoBody)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Create(c)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["todo"], "Title")

	// NOTE: Todoリストが作成されていることを確認
	todo := models.Todo{}
	if err := DBCon.Where("user_id = ?", user.ID).First(&todo).Error; err != nil {
		s.T().Fatalf("failed to create todo %v", err)
	}
	assert.Equal(s.T(), "test title 1", todo.Title)
	assert.Equal(s.T(), "test content 1", todo.Content)
}

func (s *TestTodoControllerSuite) TestCreateTodo_ValidationError() {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	createTodoBody := bytes.NewBufferString("{\"title\":\"\",\"content\":\"test content 1\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/todos", createTodoBody)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Create(c)

	assert.Equal(s.T(), 400, res.Code)

	// NOTE: Todoリストが作成されていないことを確認
	todo := models.Todo{}
	err := DBCon.Where("user_id = ?", user.ID).First(&todo).Error
	assert.NotNil(s.T(), err)
}

func (s *TestTodoControllerSuite) TestIndex() {
	// NOTE: Todoのデータを作っておく
	todos := []models.Todo{
		{Title: "test title 1", Content: "test content 1", UserID: user.ID},
		{Title: "test title 2", Content: "test content 2", UserID: user.ID},
	}
	if err := DBCon.Create(&todos).Error; err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request, _ = http.NewRequest(http.MethodGet, "/todos", nil)
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Index(c)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Len(s.T(), responseBody["todos"], 2)
}

func (s *TestTodoControllerSuite) TestShow() {
	// NOTE: Todoのデータを作っておく
	todo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DBCon.Create(&todo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	todoID := strconv.Itoa(todo.ID)
	param := gin.Param{Key: "id", Value: todoID}
	c.Params = gin.Params{param}
	c.Request, _ = http.NewRequest(http.MethodGet, "/todos/"+todoID, nil)
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Show(c)

	assert.Equal(s.T(), 200, res.Code)
}

func (s *TestTodoControllerSuite) TestUpdate() {
	// NOTE: Todoのデータを作っておく
	todo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DBCon.Create(&todo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	todoID := strconv.Itoa(todo.ID)
	param := gin.Param{Key: "id", Value: todoID}
	c.Params = gin.Params{param}
	updateTodoBody := bytes.NewBufferString("{\"title\":\"test updated title 1\",\"content\":\"test updated content 1\"}")
	c.Request, _ = http.NewRequest(http.MethodPut, "/todos/"+todoID, updateTodoBody)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Update(c)

	assert.Equal(s.T(), 200, res.Code)
	// NOTE: Todoリストが更新されていることを確認
	updatedTodo := models.Todo{}
	if err := DBCon.Where("user_id = ?", user.ID).First(&updatedTodo).Error; err != nil {
		s.T().Fatalf("failed to create todo %v", err)
	}
	assert.Equal(s.T(), "test updated title 1", updatedTodo.Title)
	assert.Equal(s.T(), "test updated content 1", updatedTodo.Content)
}

func (s *TestTodoControllerSuite) TestUpdateTodo_ValidationError() {
	// NOTE: Todoのデータを作っておく
	todo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DBCon.Create(&todo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	todoID := strconv.Itoa(todo.ID)
	param := gin.Param{Key: "id", Value: todoID}
	c.Params = gin.Params{param}
	updateTodoBody := bytes.NewBufferString("{\"title\":\"\",\"content\":\"test content 1\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/todos/"+todoID, updateTodoBody)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Create(c)

	assert.Equal(s.T(), 400, res.Code)

	// NOTE: Todoが更新されていないこと
	updatedTodo := models.Todo{}
	if err := DBCon.Where("user_id = ?", user.ID).First(&updatedTodo).Error; err != nil {
		s.T().Fatalf("failed to create todo %v", err)
	}
	assert.Equal(s.T(), "test title 1", updatedTodo.Title)
}

func (s *TestTodoControllerSuite) TestDelete() {
	// NOTE: Todoのデータを作っておく
	todo := models.Todo{Title: "test title 1", Content: "test content 1", UserID: user.ID}
	if err := DBCon.Create(&todo).Error; err != nil {
		s.T().Fatalf("failed to create test todo %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	todoID := strconv.Itoa(todo.ID)
	param := gin.Param{Key: "id", Value: todoID}
	c.Params = gin.Params{param}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/todos/"+todoID, nil)
	c.Request.Header.Set("Cookie", "token="+token)
	testTodoController.Delete(c)

	assert.Equal(s.T(), 200, res.Code)
	// NOTE: Todoリストが削除されていることを確認
	deletedTodo := models.Todo{}
	err := DBCon.Where("user_id = ?", user.ID).First(&deletedTodo).Error
	assert.NotNil(s.T(), err)
}

func TestTodoController(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestTodoControllerSuite))
}
