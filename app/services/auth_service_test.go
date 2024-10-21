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

type TestAuthServiceSuite struct {
	WithDBSuite
}

var testAuthService AuthService

func (s *TestAuthServiceSuite) SetupTest() {
	s.SetDBCon()

	userRepository := repositories.NewUserRepository(DBCon)
	testAuthService = NewAuthService(userRepository)
}

func (s *TestAuthServiceSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestAuthServiceSuite) TestSignUp() {
	requestParams := dto.SignUpRequest{Name: "test name 1", Email: "test@example.com", Password: "password"}

	result := testAuthService.SignUp(requestParams)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.ErrorType)

	// NOTE: ユーザが作成されていることを確認
	user := models.User{}
	if err := DBCon.First(&user).Error; err != nil {
		s.T().Fatalf("failed to create user %v", err)
	}
	assert.Equal(s.T(), "test name 1", user.Name)
	assert.Equal(s.T(), "test@example.com", user.Email)
}

func (s *TestAuthServiceSuite) TestSignUp_ValidationError() {
	requestParams := dto.SignUpRequest{Name: "test name 1", Email: "", Password: "password"}

	result := testAuthService.SignUp(requestParams)

	assert.NotNil(s.T(), result.Error)
	assert.Equal(s.T(), "validationError", result.ErrorType)

	// NOTE: ユーザが作成されていないことを確認
	user := models.User{}
	err := DBCon.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *TestAuthServiceSuite) TestSignIn() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DBCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	requestParams := dto.SignInRequest{Email: "test@example.com", Password: "password"}

	result := testAuthService.SignIn(requestParams)

	assert.Nil(s.T(), result.Error)
	assert.Equal(s.T(), "", result.NotFoundMessage)
	assert.NotNil(s.T(), result.TokenString)
}

func (s *TestAuthServiceSuite) TestSignIn_NotFoundError() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DBCon.Create(&user).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	requestParams := dto.SignInRequest{Email: "test_1@example.com", Password: "password"}

	result := testAuthService.SignIn(requestParams)

	assert.Equal(s.T(), "メールアドレスまたはパスワードに該当するユーザが存在しません。", result.NotFoundMessage)
}

func TestAuthService(t *testing.T) {
	// テストスイートを実行
	suite.Run(t, new(TestAuthServiceSuite))
}
