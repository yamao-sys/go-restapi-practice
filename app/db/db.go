package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	err error
)

func Init() *gorm.DB {
	// DBインスタンス生成
	Db, err = gorm.Open(mysql.Open(GetDsn()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return Db
}

func Close(db *gorm.DB) {
	sqlDb, _ := db.DB()
	if err := sqlDb.Close(); err != nil {
		panic(err)
	}
}
