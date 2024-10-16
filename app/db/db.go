package db

import (
	"app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	err error
)

func Init() *gorm.DB {
	dsn := config.Config.DbUserName +
		":" +
		config.Config.DbUserPassword +
		"@tcp(" + config.Config.DbHost + ":" + config.Config.DbPort + ")/" +
		config.Config.DbName +
		"?charset=utf8mb4&parseTime=true&loc=Local"

	// DBインスタンス生成
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
