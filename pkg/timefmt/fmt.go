package timefmt

import (
	"fmt"
	"time"
)

func ElapsedTime(since time.Time) string {
	prefix := "Elapsed time:"
	duration := time.Since(since)
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%v %2vd %2vh %2vm %2vs", prefix, days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%v %2vh %2vm %2vs", prefix, hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%v %2vm %2vs", prefix, minutes, seconds)
	}
	return fmt.Sprintf("%v %2vs", prefix, seconds)
}
