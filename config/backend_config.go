package config

import (
	"github.com/spf13/viper"
)

type BackendConfig struct {
	Database  DatabaseConfig
	JwtSecret []byte
	CorsUrls  string
}

var BackendCFG *BackendConfig

func LoadBackendConfig() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	BackendCFG = &BackendConfig{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JwtSecret: []byte("jwt_secret"),
		CorsUrls:  viper.GetString("CORS_URLS"),
	}
}
