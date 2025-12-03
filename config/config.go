package config

import (
	"log"
	"time"
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
	
	AlertReminderIntervalMinutes int64
	Alert7DaysMinutes            int64
	Alert1DayMinutes             int64
	Alert1HourMinutes            int64
	AlertScheduledTime           string
	AlertScheduledHour           int
	AlertScheduledMinute         int
}

var CFG *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	alertReminderInterval := viper.GetInt64("ALERT_REMINDER_INTERVAL_MINUTES")
	if alertReminderInterval == 0 {
		alertReminderInterval = 1440
	}
	
	alert7Days := viper.GetInt64("ALERT_7DAYS_MINUTES")
	if alert7Days == 0 {
		alert7Days = 10080
	}
	
	alert1Day := viper.GetInt64("ALERT_1DAY_MINUTES")
	if alert1Day == 0 {
		alert1Day = 1440
	}
	
	alert1Hour := viper.GetInt64("ALERT_1HOUR_MINUTES")
	if alert1Hour == 0 {
		alert1Hour = 60
	}

	var alertScheduledTime string
	var alertScheduledHour, alertScheduledMinute int
	
	if viper.IsSet("ALERT_SCHEDULED_TIME") {
		alertScheduledTime = viper.GetString("ALERT_SCHEDULED_TIME")
		parsedTime, err := time.Parse("15:04", alertScheduledTime)
		if err != nil {
			log.Printf("Warning: ALERT_SCHEDULED_TIME=%s is invalid (expected HH:MM format), using default 12:00", alertScheduledTime)
			alertScheduledTime = "12:00"
			alertScheduledHour = 12
			alertScheduledMinute = 0
		} else {
			alertScheduledHour = parsedTime.Hour()
			alertScheduledMinute = parsedTime.Minute()
		}
	} else {
		alertScheduledTime = "12:00"
		alertScheduledHour = 12
		alertScheduledMinute = 0
	}

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
		AlertReminderIntervalMinutes: alertReminderInterval,
		Alert7DaysMinutes:            alert7Days,
		Alert1DayMinutes:             alert1Day,
		Alert1HourMinutes:            alert1Hour,
		AlertScheduledTime:           alertScheduledTime,
		AlertScheduledHour:           alertScheduledHour,
		AlertScheduledMinute:         alertScheduledMinute,
	}
}
