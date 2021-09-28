package util

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	DBHost              string
	DBDriver            string
	DBUser              string
	DBPassword          string
	DBName              string
	DBPort              string
	BinPath             string
	ViewPath            string
	BaseUrl             string
	AccessTokenDuration time.Duration
	TokenSymmetricKey   string
	FeUrl               string
}

func LoadConfig(path string) (config Config, err error) {
	err = godotenv.Load(filepath.Join(path, ".env"))

	if err != nil {
		return
	}

	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DBHost = os.Getenv("DB_HOST")
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")
	config.DBPort = os.Getenv("DB_PORT")
	config.BinPath = os.Getenv("BIN_PATH")
	config.ViewPath = os.Getenv("VIEW_PATH")
	config.BaseUrl = os.Getenv("BASE_URL")
	config.AccessTokenDuration, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	config.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	config.FeUrl = os.Getenv("FE_URL")

	return
}
