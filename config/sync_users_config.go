package config

import (
	"github.com/spf13/viper"
)

type SyncUsersConfig struct {
	Database DatabaseConfig
	BotToken string
}

var SyncUsersCFG *SyncUsersConfig

func LoadSyncUsersConfig() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	SyncUsersCFG = &SyncUsersConfig{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		BotToken: viper.GetString("TELEGRAM_BOT_TOKEN"),
	}
}
