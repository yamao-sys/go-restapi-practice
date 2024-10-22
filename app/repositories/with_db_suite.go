package repositories

import (
	"app/db"
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-txdb"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WithDBSuite struct {
	suite.Suite
}

var DBCon *gorm.DB

// func (s *WithDBSuite) SetupSuite()                           {} // テストスイート実施前の処理
// func (s *WithDBSuite) TearDownSuite()                        {} // テストスイート終了後の処理
// func (s *WithDBSuite) SetupTest()                            {} // テストケース実施前の処理
// func (s *WithDBSuite) TearDownTest()                         {} // テストケース終了後の処理
// func (s *WithDBSuite) BeforeTest(suiteName, testName string) {} // テストケース実施前の処理
// func (s *WithDBSuite) AfterTest(suiteName, testName string)  {} // テストケース終了後の処理

func init() {
	txdb.Register("txdb", "mysql", db.GetDsn())
}

func (s *WithDBSuite) SetDBCon() {
	db, err := sql.Open("txdb", "connect")
	if err != nil {
		log.Fatalln(err)
	}

	DBCon, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to initialize GORM DB: %v", err)
	}
}

func (s *WithDBSuite) CloseDB() {
	db, _ := DBCon.DB()
	db.Close()
}
