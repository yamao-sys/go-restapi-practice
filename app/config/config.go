package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigList struct {
	DbDriverName   string
	DbName         string
	DbUserName     string
	DbUserPassword string
	DbHost         string
	DbPort         string
	ServerPort     int
}

var Config ConfigList

func init() {
	var envFilePath string
	if os.Getenv("ENV") != "" {
		envFilePath = "/app/.env." + os.Getenv("ENV")
	} else {
		envFilePath = "/app/.env.development"
	}
	godotenv.Load(envFilePath)

	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	Config = ConfigList{
		DbDriverName:   os.Getenv("DB_DRIVER_NAME"),
		DbName:         os.Getenv("DB_NAME"),
		DbUserName:     os.Getenv("DB_USER_NAME"),
		DbUserPassword: os.Getenv("DB_USER_PASSWORD"),
		DbHost:         os.Getenv("DB_HOST"),
		DbPort:         os.Getenv("DB_PORT"),
		ServerPort:     serverPort,
	}
}
