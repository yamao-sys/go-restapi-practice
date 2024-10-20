package controllers

import (
	"app/config"
	"app/models"
	"app/repositories"
	"app/services"
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WithDbSuite struct {
	suite.Suite
}

var (
	DbCon *gorm.DB
	pid   int
	token string
)

// func (s *WithDbSuite) SetupSuite()                           {} // テストスイート実施前の処理
// func (s *WithDbSuite) TearDownSuite()                        {} // テストスイート終了後の処理
// func (s *WithDbSuite) SetupTest()                            {} // テストケース実施前の処理
// func (s *WithDbSuite) TearDownTest()                         {} // テストケース終了後の処理
// func (s *WithDbSuite) BeforeTest(suiteName, testName string) {} // テストケース実施前の処理
// func (s *WithDbSuite) AfterTest(suiteName, testName string)  {} // テストケース終了後の処理

func init() {
	pid = os.Getpid()

	dsn := config.Config.DbUserName +
		":" +
		config.Config.DbUserPassword +
		"@tcp(" + config.Config.DbHost + ":" + config.Config.DbPort + ")/" +
		"go_restapi_practice_test" +
		"?charset=utf8mb4&parseTime=true&loc=Local"

	txdb.Register("txdb-controller", "mysql", dsn)
}

func (s *WithDbSuite) SetDbCon() {
	log.Printf("pid: %v", pid)
	db, err := sql.Open("txdb-controller", "connect"+strconv.Itoa(pid))
	if err != nil {
		log.Fatalln(err)
	}

	DbCon, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to initialize GORM DB: %v", err)
	}
	DbCon.AutoMigrate(&models.User{}, &models.Todo{})
}

func (s *WithDbSuite) CloseDb() {
	db, _ := DbCon.DB()
	db.Close()
}

func (s *WithDbSuite) signIn() {
	userRepository := repositories.NewUserRepository(DbCon)
	authService := services.NewAuthService(userRepository)
	authController := NewAuthController(authService)

	// gin contextの生成
	authRecorder := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(authRecorder)

	// NOTE: リクエストの生成
	body := bytes.NewBufferString("{\"email\":\"test@example.com\",\"password\":\"password\"}")
	req, _ := http.NewRequest("POST", "/auth/sign_in", body)
	req.Header.Set("Content-Type", "application/json")
	ginContext.Request = req

	// NOTE: ログインし、tokenに認証情報を格納
	authController.SignIn(ginContext)
	token = authRecorder.Result().Cookies()[0].Value
}
