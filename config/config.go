package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DbURL  string
	Debug  bool
	Volume string
	Port   string
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Info("No .env getting from actual env")
	}

	conf := Config{
		DbURL:  os.Getenv("DB_URL"),
		Debug:  (os.Getenv("Debug") == "true"),
		Volume: os.Getenv("VOLUME"),
		Port:   os.Getenv("PORT"),
	}

	if conf.Volume == "" || conf.DbURL == "" || conf.Port == "" {
		logrus.Fatalf("One or more of the required env variables are not set conf: %+v", conf)
	}

	return conf
}
