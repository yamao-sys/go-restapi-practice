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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	testAuthController AuthController
)

type TestAuthControllerSuite struct {
	WithDBSuite
}

func (s *TestAuthControllerSuite) SetupTest() {
	s.SetDBCon()

	userRepository := repositories.NewUserRepository(DBCon)

	authService := services.NewAuthService(userRepository)

	// NOTE: テスト対象のコントローラを設定
	testAuthController = NewAuthController(authService)
}

func (s *TestAuthControllerSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestAuthControllerSuite) TestSignUp() {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	signUpRequestBody := bytes.NewBufferString("{\"name\":\"test name 1\",\"email\":\"test@example.com\",\"password\":\"password\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/sign_up", signUpRequestBody)
	c.Request.Header.Set("Content-Type", "application/json")
	testAuthController.SignUp(c)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["user"], "Name")

	// NOTE: ユーザが作成されていることを確認
	user := models.User{}
	if err := DBCon.Where("email = ?", "test@example.com").First(&user).Error; err != nil {
		s.T().Fatalf("failed to create todo %v", err)
	}
	assert.Equal(s.T(), "test name 1", user.Name)
}

func (s *TestAuthControllerSuite) TestSignUp_ValidationError() {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	signUpRequestBody := bytes.NewBufferString("{\"name\":\"test name 1\",\"email\":\"\",\"password\":\"password\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/sign_up", signUpRequestBody)
	c.Request.Header.Set("Content-Type", "application/json")
	testAuthController.SignUp(c)

	assert.Equal(s.T(), 400, res.Code)

	// NOTE: ユーザが作成されていないことを確認
	user := models.User{}
	err := DBCon.Where("email = ?", "test@example.com").First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *TestAuthControllerSuite) TestSignIn() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DBCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	signInRequestBody := bytes.NewBufferString("{\"email\":\"test@example.com\",\"password\":\"password\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/sign_in", signInRequestBody)
	c.Request.Header.Set("Content-Type", "application/json")
	testAuthController.SignIn(c)

	assert.Equal(s.T(), 200, res.Code)
	token = res.Result().Cookies()[0].Value
	assert.NotEmpty(s.T(), token)
}

func (s *TestAuthControllerSuite) TestSignIn_NotFoundError() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DBCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	signInRequestBody := bytes.NewBufferString("{\"email\":\"test_1@example.com\",\"password\":\"password\"}")
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/sign_in", signInRequestBody)
	c.Request.Header.Set("Content-Type", "application/json")
	testAuthController.SignIn(c)

	assert.Equal(s.T(), 404, res.Code)
	assert.Empty(s.T(), res.Result().Cookies())
}

func TestAuthController(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestAuthControllerSuite))
}
