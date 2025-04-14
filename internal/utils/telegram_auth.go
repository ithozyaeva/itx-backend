package utils

import "time"

func CheckExpirationDate(date time.Time) bool {
	return date.After(time.Now().AddDate(0, 1, 0))
}
