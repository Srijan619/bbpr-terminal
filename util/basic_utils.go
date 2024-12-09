package util

import (
	"simple-git-terminal/types"
	"strings"
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
