package utils

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

func ParseDate(dateStr *string) (*DateOnly, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse(dateFormat, *dateStr)
	if err != nil {
		return nil, errors.New("invalid date format (expected YYYY-MM-DD)")
	}
	return NewDateOnly(&parsedDate), nil
}

func NewDateOnly(t *time.Time) *DateOnly {
	if t == nil {
		return nil
	}
	d := DateOnly(*t)
	return &d
}

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + t.Format(dateFormat) + `"`), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		return nil
	}

	t, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format(dateFormat), nil
}

func (d *DateOnly) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = DateOnly(v)
		return nil
	case []byte:
		t, err := time.Parse(dateFormat, string(v))
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	case string:
		t, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	default:
		return fmt.Errorf("cannot scan value %v into DateOnly", value)
	}
}
