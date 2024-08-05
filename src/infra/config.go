package infra

import "os"

var osGetenv = os.Getenv

type SqlDbConf struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     string
	SSLMode  string
	Schema   string
}

type Config struct {
	App   string
	SqlDb SqlDbConf
}

func Make() Config {
	sqldb := SqlDbConf{
		Host:     osGetenv("DB_HOST"),
		Username: osGetenv("DB_USERNAME"),
		Password: osGetenv("DB_PASSWORD"),
		Name:     osGetenv("DB_NAME"),
		Port:     osGetenv("DB_PORT"),
		SSLMode:  osGetenv("DB_SSL_MODE"),
		Schema:   osGetenv("DB_SCHEMA"),
	}

	return Config{
		App:   osGetenv("APP_NAME"),
		SqlDb: sqldb,
	}
}
