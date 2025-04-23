package utils

import (
	"fmt"
	"time"
)

func ParseDate(dateStr *string) (*time.Time, error) {
	if dateStr != nil && *dateStr == "" {
		parsedDate, err := time.Parse("2006-01-02", *dateStr)
		if err != nil {
			return nil, fmt.Errorf("НЕВЕРНЫЙ ФОРМАТ ДАТЫ. ИСПОЛЬЗУЙТЕ ФОРМАТ YYYY-MM-DD")
		}
		return &parsedDate, nil
	}
	return nil, nil
}
