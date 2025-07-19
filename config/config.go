package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port       string
	Host       string
	CHHost     string
	CHPort     string
	CHUser     string
	CHPassword string
	CHHandler  string
	CSHost     string
	CSPort     string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		return nil, err
	}
	return &Config{
		Port:       os.Getenv("PORT"),
		Host:       os.Getenv("HOST"),
		CHHost:     os.Getenv("CH_HOST"),
		CHPort:     os.Getenv("CH_PORT"),
		CHUser:     os.Getenv("CH_USER"),
		CHPassword: os.Getenv("CH_PASSWORD"),
		CHHandler:  os.Getenv("CH_HANDLER"),
		CSHost:     os.Getenv("CS_HOST"),
		CSPort:     os.Getenv("CS_PORT"),
	}, err
}
