package db

import (
	"app/config"
)

func GetDsn() string {
	return config.Config.DBUserName +
		":" +
		config.Config.DBUserPassword +
		"@tcp(" + config.Config.DBHost + ":" + config.Config.DBPort + ")/" +
		config.Config.DBName +
		"?charset=utf8mb4&parseTime=true&loc=Local"
}
