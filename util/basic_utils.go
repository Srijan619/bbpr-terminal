package util

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"simple-git-terminal/types"
	"strings"
	"time"
	"unicode"
)

func removeZeroWidth(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	for _, r := range input {
		if !unicode.Is(unicode.Mn, r) && r != '\u200C' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func SanitizePR(pr types.PR) types.PR {
	pr.Title = removeZeroWidth(pr.Title)
	if desc, ok := pr.Description.(string); ok {
		pr.Description = removeZeroWidth(desc)
	}
	return pr
}

// Helper function to format initials with a distinct color
func FormatInitials(initials string) string {
	return fmt.Sprintf("[::b]%s[-]", getInitials(initials))
}

// Get the initials of the author's display name
func getInitials(displayName string) string {
	words := strings.Fields(displayName)
	if len(words) > 0 {
		initials := ""
		for _, word := range words {
			initials += string(word[0])
		}
		return strings.ToUpper(initials)
	}

	if len(displayName) > 1 {
		return strings.ToUpper(displayName[:2])
	}
	return strings.ToUpper(displayName)
}

// Helper function to calculate time ago
func FormatTimeAgo(date string) string {
	parsedTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "unknown time"
	}
	return humanize.Time(parsedTime)
}

func FormatCombinedTimeAgo(date string) string {
	parsedTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "unknown time"
	}
	formattedDate := parsedTime.Format("2006-01-02")
	return fmt.Sprintf("%s (%s)", formattedDate, FormatTimeAgo(date))
}
