package utils

import (
	"fmt"
	"ithozyeva/internal/models"
	"strings"
	"time"
)

func GenerateICS(event *models.Event) string {
	formatTime := func(t time.Time) string {
		return t.UTC().Format("20060102T150405Z") // iCalendar формат
	}

	builder := strings.Builder{}
	builder.WriteString("BEGIN:VCALENDAR\n")
	builder.WriteString("VERSION:2.0\n")

	builder.WriteString("BEGIN:VEVENT\n")
	builder.WriteString(fmt.Sprintf("UID:event-%d@example.com\n", event.Id))
	builder.WriteString(fmt.Sprintf("DTSTAMP:%s\n", formatTime(time.Now())))
	builder.WriteString(fmt.Sprintf("DTSTART:%s\n", formatTime(event.Date)))
	builder.WriteString(fmt.Sprintf("SUMMARY:%s\n", escapeICS(event.Title)))
	builder.WriteString(fmt.Sprintf("DESCRIPTION:%s\n", escapeICS(event.Description)))

	// Место проведения
	place := event.Place
	if event.PlaceType == models.EventHybrid && event.CustomPlaceType != "" {
		place = event.CustomPlaceType + ": " + event.Place
	}
	// Видеоссылка, если есть
	if event.PlaceType == models.EventOnline {
		builder.WriteString(fmt.Sprintf("LOCATION:%s\n", escapeICS(place)))
		builder.WriteString(fmt.Sprintf("URL:%s\n", escapeICS(event.Place)))
	} else {
		builder.WriteString(fmt.Sprintf("LOCATION:%s\n", escapeICS(place)))
	}

	builder.WriteString("END:VEVENT\n")
	builder.WriteString("END:VCALENDAR\n")
	return builder.String()
}
func escapeICS(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		";", "\\;",
		",", "\\,",
		"\n", "\\n",
	)
	return replacer.Replace(s)
}
