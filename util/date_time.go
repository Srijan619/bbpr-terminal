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
