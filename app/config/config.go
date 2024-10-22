package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigList struct {
	DBDriverName   string
	DBName         string
	DBUserName     string
	DBUserPassword string
	DBHost         string
	DBPort         string
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
		DBDriverName:   os.Getenv("DB_DRIVER_NAME"),
		DBName:         os.Getenv("DB_NAME"),
		DBUserName:     os.Getenv("DB_USER_NAME"),
		DBUserPassword: os.Getenv("DB_USER_PASSWORD"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		ServerPort:     serverPort,
	}
}
