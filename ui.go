package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	HIGH_CONTRAST_COLOR = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR  = tcell.ColorYellow
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
		SetColumns(60, 0).
		SetBorders(true)

	mainGrid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	mainGrid.AddItem(prList, 1, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(prDetails, 1, 1, 1, 1, 0, 0, false)

	// Populate PR list
	populatePRList(prs, prList)

	// Handle PR selection
	prList.SetSelectedFunc(func(row, column int) {
		updatePRDetails(prs, prDetails, row)
	})

	// Set initial PR details
	if len(prs) > 0 {
		prList.Select(0, 0)
		updatePRDetails(prs, prDetails, 0)
	}

	app.SetRoot(mainGrid, true)
	return app
}

func updatePRDetails(prs []PR, prDetails *tview.TextView, row int) {
	if row >= 0 && row < len(prs) {
		selectedPR := prs[row]
		description := formatDescription(selectedPR.Description)

		// Build the PR details
		prDetails.SetText(fmt.Sprintf(
			"[::b]Title:[-] %s\n"+
				"[::b]State:[-] %s\n"+
				"[::b]Author:[-] %s\n"+
				"[::b]Created On:[-] %s\n"+
				"[::b]Updated On:[-] %s\n"+
				"[::b]Link:[-] %s\n"+
				"[::b]Description:[-] %s\n",
			selectedPR.Title,
			selectedPR.State,
			selectedPR.Author.DisplayName,
			selectedPR.CreatedOn,
			selectedPR.UpdatedOn,
			selectedPR.Links.HTML.Href,
			description,
		))
	}
}

// Function to populate the PR list
func populatePRList(prs []PR, prList *tview.Table) {
	for i, pr := range prs {
		titleCell := cellFormat(ellipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := styleState(pr.State)

		initialsCell := cellFormat(formatInitials(pr.Author.DisplayName), HIGH_CONTRAST_COLOR)

		sourceBranch := cellFormat(ellipsizeText(pr.Source.Branch.Name, 10), LOW_CONTRAST_COLOR)
		arrow := cellFormat("-->", LOW_CONTRAST_COLOR)
		destinationBranch := cellFormat(ellipsizeText(pr.Destination.Branch.Name, 10), LOW_CONTRAST_COLOR)

		prList.SetCell(i, 0, initialsCell)
		prList.SetCell(i, 1, stateCell)
		prList.SetCell(i, 2, titleCell)

		prList.SetCell(i, 3, sourceBranch)
		prList.SetCell(i, 4, arrow)
		prList.SetCell(i, 5, destinationBranch)
	}
}

func cellFormat(text string, color tcell.Color) *tview.TableCell {
	return tview.NewTableCell(text).
		SetTextColor(color).
		SetAlign(tview.AlignLeft).
		SetSelectable(true)
}

// Helper function to format initials with a distinct color
func formatInitials(initials string) string {
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

// Helper function to style PR state
func styleState(state string) *tview.TableCell {
	var stateColor tcell.Color

	switch state {
	case "OPEN":
		stateColor = tcell.ColorGreen
	case "MERGED":
		stateColor = tcell.ColorBlue
	case "DECLINED":
		stateColor = tcell.ColorRed
	default:
		stateColor = tcell.ColorYellow
	}

	return tview.NewTableCell(state).
		SetTextColor(stateColor).
		SetAlign(tview.AlignLeft).
		SetSelectable(true)
}

func ellipsizeText(text string, max int) string {
	if len(text) > max {
		text = text[:max] + "..."
	}
	return text
}
