package config

import (
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	Database      DatabaseConfig
	JwtSecret     []byte
	CorsUrls      string
	Port          string
	DatabaseURL   string
	TelegramToken string
}

var CFG *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	CFG = &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JwtSecret:     []byte("jwt_secret"),
		CorsUrls:      viper.GetString("CORS_URLS"),
		Port:          viper.GetString("PORT"),
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		TelegramToken: viper.GetString("TELEGRAM_BOT_TOKEN"),
	}
}
