package db

import (
	"app/config"
)

func GetDsn() string {
	return config.Config.DbUserName +
		":" +
		config.Config.DbUserPassword +
		"@tcp(" + config.Config.DbHost + ":" + config.Config.DbPort + ")/" +
		config.Config.DbName +
		"?charset=utf8mb4&parseTime=true&loc=Local"
}
