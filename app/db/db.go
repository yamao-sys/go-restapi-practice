package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() *gorm.DB {
	// DBインスタンス生成
	DB, err = gorm.Open(mysql.Open(GetDsn()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return DB
}

func Close(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}
