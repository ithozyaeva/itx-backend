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
	Database           DatabaseConfig
	JwtSecret          []byte
	Port               string
	TelegramToken      string
	TelegramMainChatID int64
	PublicDomain       string
	BackendDomain      string
	S3                 S3Config
}

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
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
		JwtSecret:          []byte("jwt_secret"),
		Port:               viper.GetString("PORT"),
		TelegramToken:      viper.GetString("TELEGRAM_BOT_TOKEN"),
		TelegramMainChatID: viper.GetInt64("TELEGRAM_MAIN_CHAT_ID"),
		PublicDomain:       viper.GetString("PUBLIC_DOMAIN"),
		BackendDomain:      viper.GetString("BACKEND_DOMAIN"),
		S3: S3Config{
			Endpoint:  viper.GetString("S3_ENDPOINT"),
			Region:    viper.GetString("S3_REGION"),
			AccessKey: viper.GetString("S3_ACCESS_KEY"),
			SecretKey: viper.GetString("S3_SECRET_KEY"),
			Bucket:    viper.GetString("S3_BUCKET"),
		},
	}
}
