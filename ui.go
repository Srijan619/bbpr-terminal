package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"simple-git-terminal/pr"
	"simple-git-terminal/types"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	HIGH_CONTRAST_COLOR = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR  = tcell.ColorYellow
)

func CreateApp(prs []types.PR) *tview.Application {
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
		SetDynamicColors(true)

	activityDetails := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Grid layout
	mainGrid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(60, 0).
		SetBorders(true)

	mainGrid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	mainGrid.AddItem(prList, 1, 0, 1, 1, 0, 0, true)

	rightPanelGrid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(60, 0).
		SetBorders(true)

	rightPanelHeader := tview.NewTextView().SetText("Here we will display PR title").SetTextAlign(tview.AlignCenter)

	rightPanelGrid.AddItem(rightPanelHeader, 0, 0, 1, 1, 0, 0, false)
	rightPanelGrid.AddItem(prDetails, 1, 0, 1, 1, 0, 0, false)
	rightPanelGrid.AddItem(activityDetails, 2, 0, 1, 1, 0, 0, false)

	mainGrid.AddItem(rightPanelGrid, 1, 1, 1, 1, 0, 0, false)
	// Populate PR list
	populatePRList(prs, prList)

	prList.SetSelectedFunc(func(row, column int) {
		pr.UpdatePrDetails(prs, prDetails, row)

		if row >= 0 && row < len(prs) {
			selectedPR := prs[row]
			go func() {
				rightPanelHeader.SetText(selectedPR.Title)

				activityDetails.Clear()
				activityDetails.AddItem(tview.NewTextView().SetText("⏳ Fetching activities..."), 0, 1, true)

				prActivities := fetchBitbucketActivities(selectedPR.ID)

				app.QueueUpdateDraw(func() {
					activityDetails.Clear()
					activityDetails.AddItem(pr.CreateActivitiesView(prActivities), 0, 1, true)
				})
			}()
		}
	})

	// Set initial PR details
	if len(prs) > 0 {
		prList.Select(0, 0)
		pr.UpdatePrDetails(prs, prDetails, 0)

		// Fetch initial activities dynamically
		initialPR := prs[0]
		go func() {

			rightPanelHeader.SetText(initialPR.Title)
			prActivities := fetchBitbucketActivities(initialPR.ID)
			app.QueueUpdateDraw(func() {
				activityDetails.AddItem(pr.CreateActivitiesView(prActivities), 0, 1, true)
			})
		}()
	}

	app.SetRoot(mainGrid, true)
	return app
}

// Function to populate the PR list
func populatePRList(prs []types.PR, prList *tview.Table) {
	for i, pr := range prs {
		titleCell := cellFormat(util.EllipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := util.CreateStateCell(pr.State)

		initialsCell := cellFormat(formatInitials(pr.Author.DisplayName), HIGH_CONTRAST_COLOR)

		sourceBranch := cellFormat(util.EllipsizeText(pr.Source.Branch.Name, 10), LOW_CONTRAST_COLOR)
		arrow := cellFormat("-->", LOW_CONTRAST_COLOR)
		destinationBranch := cellFormat(util.EllipsizeText(pr.Destination.Branch.Name, 10), LOW_CONTRAST_COLOR)

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

// GetStateColor determines the color based on the PR state
