package util

import (
	"log"
	"net/url"
)

func ExtractQueryFromNextURL(nextURL string) string {
	parsed, err := url.Parse(nextURL)
	if err != nil {
		log.Printf("[ERROR] Failed to parse next URL: %v", err)
		return ""
	}
	return parsed.RawQuery // returns "page=2&pagelen=10" etc.
}
