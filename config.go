package main

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) ConnectionInfo() string {
	psqlInfo := os.Getenv("DATABASE_URL")
	if len(psqlInfo) == 0 {
		psqlInfo = fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", c.User, c.Password, c.Host, c.Port, c.Name)
	}
	return psqlInfo
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "1234",
		Name:     "simplegallery_dev",
	}
}

type Config struct {
	Port string
	Env  string
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func (c Config) GetPort() string {
	serverPort := os.Getenv("PORT")
	if len(serverPort) == 0 {
		serverPort = c.Port
	}
	return serverPort
}

func DefaultConfig() Config {
	return Config{
		Port: "3333",
		Env:  "dev",
	}
}
