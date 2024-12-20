package util

import (
	"bytes"
	"github.com/charmbracelet/glamour"
	"github.com/rivo/tview"
	"log"
	"strings"
)

// Global variable to store the renderer instance
var renderer *glamour.TermRenderer

// Initialize the renderer once and reuse it
func InitMdRenderer() {
	var err error
	renderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		log.Fatalf("Error initializing renderer: %v", err)
	}
}

func RenderMarkdown(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if renderer == nil {
		InitMdRenderer()
	}
	out, err := renderer.Render(trimmed)
	log.Printf("Render markdown output: %s ", out)

	if err != nil {
		log.Fatalf("Error rendering markdown: %v", err)
	}
	return strings.TrimSpace(TranslateANSI(out))
}

// Translate ANSI escape sequences into tview-compatible format
func TranslateANSI(input string) string {
	var buf bytes.Buffer
	w := tview.ANSIWriter(&buf)
	_, err := w.Write([]byte(input))
	if err != nil {
		log.Fatalf("Error translating ANSI: %v", err)
	}
	return buf.String()
}
