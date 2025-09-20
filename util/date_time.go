package util

import (
	"fmt"
	"time"
)

func HumanizeTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	default:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}

func FormatTime(input string) string {
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return input
	}
	return t.Format("2006-01-02 15:04:05")
}
