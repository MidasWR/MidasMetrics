package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port       string
	Host       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBHandler  string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		return nil, err
	}
	return &Config{
		Port:       os.Getenv("PORT"),
		Host:       os.Getenv("HOST"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHandler:  os.Getenv("DB_HANDLER"),
	}, err
}
