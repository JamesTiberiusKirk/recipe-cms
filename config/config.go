package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type Config struct {
	DbURL string
	Debug bool
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env getting from actual env")
	}

	return Config{
		DbURL: os.Getenv("DB_URL"),
		Debug: (os.Getenv("Debug") == "true"),
	}
}
