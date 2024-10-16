package config

import (
	"log"
	"os"

	"gopkg.in/go-ini/ini.v1"
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
	cfg, err := ini.Load("/app/config.ini")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	Config = ConfigList{
		DbDriverName:   cfg.Section("db").Key("db_driver_name").String(),
		DbName:         cfg.Section("db").Key("db_name").String(),
		DbUserName:     cfg.Section("db").Key("db_user_name").String(),
		DbUserPassword: cfg.Section("db").Key("db_user_password").String(),
		DbHost:         cfg.Section("db").Key("db_host").String(),
		DbPort:         cfg.Section("db").Key("db_port").String(),
		ServerPort:     cfg.Section("api").Key("server_port").MustInt(),
	}
}
