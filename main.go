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

	//	prs := fetchBitbucketPRs()

	//	app := CreateApp(prs, workspace, repoSlug)
	//
	textView := tview.NewTextView()

	textView.SetDynamicColors(true).
		SetText("PR lists").
		SetTitle("PRs").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorYellow).
		SetBorderAttributes(tcell.AttrDim | tcell.AttrItalic)

	textView2 := tview.NewTextView()

	textView2.SetDynamicColors(true).
		SetText("PR Details here").
		SetTitle("Details").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorYellow).
		SetBorderAttributes(tcell.AttrDim | tcell.AttrItalic)

	textView3 := tview.NewTextView()

	textView3.SetDynamicColors(true).
		SetText("Git Diffs here").
		SetTitle("Git diffs").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorYellow).
		SetBorderAttributes(tcell.AttrDim | tcell.AttrItalic)

	flexView := tview.NewFlex()

	flexView.AddItem(textView, 0, 1, true).
		AddItem(textView2, 0, 1, true).
		AddItem(textView3, 0, 2, true)

	app := tview.NewApplication().SetRoot(flexView, true)
	app := CreateApp()
	//app := tview.NewApplication().SetRoot(pr.GenerateDiffStatTree(pr.STATIC_DATA), true)
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	log.Printf("Application ended")
}
