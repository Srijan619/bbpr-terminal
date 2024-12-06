package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


func CreateApp(prs []PR) *tview.Application {
	app := tview.NewApplication()

  // Get the current Git directory
	gitDir := getGitRepoName()

	// UI components
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[::b]Bitbucket PR Viewer - %s", gitDir))

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)
	prDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select a PR to view details").
		SetWrap(true)

	// Grid layout
	mainGrid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(40, 0).
		SetBorders(true)

	mainGrid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	mainGrid.AddItem(prList, 1, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(prDetails, 1, 1, 1, 1, 0, 0, false)

	// Populate PR list
	for i, pr := range prs {
		prRow := fmt.Sprintf("%s [%s] by %s", pr.Title, pr.State, pr.Author.DisplayName)
		cell := tview.NewTableCell(prRow).
			SetTextColor(tcell.ColorWhite).
			SetSelectable(true)

		prList.SetCell(i, 0, cell)
	}

	// Define the function to update PR details
	updatePRDetails := func(row int) {
		if row >= 0 && row < len(prs) {
			selectedPR := prs[row]
			description := formatDescription(selectedPR.Description)

			prDetails.SetText(fmt.Sprintf(
				"[::b]Title:[-] %s\n[::b]State:[-] %s\n[::b]Author:[-] %s\n[::b]Created On:[-] %s\n[::b]Updated On:[-] %s\n[::b]Description:[-] %s\n[::b]Link:[-] %s",
				selectedPR.Title,
				selectedPR.State,
				selectedPR.Author.DisplayName,
				selectedPR.CreatedOn,
				selectedPR.UpdatedOn,
				description,
				selectedPR.Links.HTML.Href,
			))
		}
	}

	// Handle PR selection
	prList.SetSelectedFunc(func(row, column int) {
		updatePRDetails(row)
	})

	// Set initial PR details
	if len(prs) > 0 {
		prList.Select(0, 0)
		updatePRDetails(0) // Manually invoke the logic for the first PR
	}

	app.SetRoot(mainGrid, true)
	return app
}



// Formats the PR description for display
func formatDescription(description interface{}) string {
	if description == nil {
		return "No description provided."
	}
	if desc, ok := description.(string); ok {
		return strings.TrimSpace(desc)
	}
	return "Unsupported description format."
}

// Get the name of the current Git repository
func getGitRepoName() string {
	// Run the 'git rev-parse --show-toplevel' command to get the root directory of the Git repo
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		// If there's an error, fallback to the current working directory name
		cwd, _ := os.Getwd()
		return fmt.Sprintf("Unknown Repository (%s)", cwd)
	}

	// Extract the repository name from the path
	repoPath := strings.TrimSpace(string(output))
	return strings.TrimSuffix(repoPath[strings.LastIndex(repoPath, "/")+1:], "\n")
}
