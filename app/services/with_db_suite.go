package services

import (
	"app/db"
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-txdb"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WithDbSuite struct {
	suite.Suite
}

var DbCon *gorm.DB

// func (s *WithDbSuite) SetupSuite()                           {} // テストスイート実施前の処理
// func (s *WithDbSuite) TearDownSuite()                        {} // テストスイート終了後の処理
// func (s *WithDbSuite) SetupTest()                            {} // テストケース実施前の処理
// func (s *WithDbSuite) TearDownTest()                         {} // テストケース終了後の処理
// func (s *WithDbSuite) BeforeTest(suiteName, testName string) {} // テストケース実施前の処理
// func (s *WithDbSuite) AfterTest(suiteName, testName string)  {} // テストケース終了後の処理

func init() {
	txdb.Register("txdb-service", "mysql", db.GetDsn())
}

func (s *WithDbSuite) SetDbCon() {
	db, err := sql.Open("txdb-service", "connect")
	if err != nil {
		log.Fatalln(err)
	}

	DbCon, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to initialize GORM DB: %v", err)
	}
}

func (s *WithDbSuite) CloseDb() {
	db, _ := DbCon.DB()
	db.Close()
}
