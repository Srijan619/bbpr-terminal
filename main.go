package main

import (
	"fmt"
	"log"
	"os"
)

var (
	workspace string
	repoSlug  string
)

func main() {
	// Open or create the log file
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)
	// Optionally add log flags (e.g., timestamp)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Log a test message to verify
	log.Printf("Application started")

	app := CreateApp()

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	log.Printf("Application ended")
}
