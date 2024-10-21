package repositories

import (
	"app/models"
	"app/test/factories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type TestUserRePositorySuite struct {
	WithDbSuite
}

func (s *TestUserRePositorySuite) SetupTest() {
	s.SetDbCon()
}

func (s *TestUserRePositorySuite) TearDownTest() {
	s.CloseDb()
}

func (s *TestUserRePositorySuite) TestCreateUser() {
	user := models.User{}
	user.Name = "test user 1"
	user.Email = "test@example.com"
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		s.T().Fatalf("failed to generate hash %v", err)
	}
	user.Password = string(hash)

	ur := NewUserRepository(DbCon)
	ur.CreateUser(&user)

	assert.NotEqual(s.T(), 0, user.ID)
}

func (s *TestUserRePositorySuite) TestFindUserByEmail() {
	testUser := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DbCon.Create(&testUser).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	user := models.User{}
	ur := NewUserRepository(DbCon)
	ur.FindUserByEmail(&user, "test@example.com")

	assert.Equal(s.T(), testUser.ID, user.ID)
}

func (s *TestUserRePositorySuite) TestFindUserById() {
	testUser := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := DbCon.Create(&testUser).Error; err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	ur := NewUserRepository(DbCon)
	user := ur.FindUserById(testUser.ID)

	assert.Equal(s.T(), testUser.Name, user.Name)
}

func TestUserRepository(t *testing.T) {
	// テストスイートを実行
	suite.Run(t, new(TestUserRePositorySuite))
}
