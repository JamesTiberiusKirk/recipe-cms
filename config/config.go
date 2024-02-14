package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DbURL    string
	Debug    bool
	Volume   string
	HTTPPort string
	Host     string
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Info("No .env getting from actual env")
	}

	conf := Config{
		DbURL:    os.Getenv("DB_URL"),
		Debug:    (os.Getenv("Debug") == "true"),
		Volume:   os.Getenv("VOLUME"),
		HTTPPort: os.Getenv("HTTP_PORT"),
		Host:     os.Getenv("HOST"),
	}

	if conf.Volume == "" ||
		conf.DbURL == "" ||
		conf.HTTPPort == "" ||
		conf.Host == "" {
		logrus.Fatalf("One or more of the required env variables are not set conf: %+v", conf)
	}

	return conf
}
