package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func (c Config) GetPort() string {
	serverPort := os.Getenv("PORT")
	if len(serverPort) == 0 {
		serverPort = strconv.Itoa(c.Port)
	}
	return serverPort
}

func DefaultConfig() Config {
	return Config{
		Port:     3333,
		Env:      "dev",
		Pepper:   "some-secret",
		HMACKey:  "hmac-secret",
		Database: DefaultPostgresConfig(),
	}
}

func LoadConfig(configReq bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if configReq {
			panic(err)
		}
		fmt.Println("Using default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .config")
	return c
}
