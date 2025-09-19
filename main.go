package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rivo/tview"
)

var (
	mode      string
	workspace string
	repoSlug  string
)

func init() {
	flag.StringVar(&mode, "mode", "pr", "Mode of the app: 'pipeline' or 'pr'")
}

func main() {
	// Parse flags
	flag.Parse()

	// Open or create the log file
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Application started in mode: %s", mode)

	var app *tview.Application

	switch mode {
	case "pr":
		app = CreateMainApp()
	case "pipeline":
		app = CreateMainAppForBBPipeline()
	default:
		log.Fatalf("Unknown mode: %s. Use 'pipeline' or 'pr'", mode)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	log.Printf("Application ended")
}
