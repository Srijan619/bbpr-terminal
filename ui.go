package main

import (
	"fmt"

	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"simple-git-terminal/pr"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

const (
	HIGH_CONTRAST_COLOR      = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR       = tcell.ColorYellow
	VIEW_ACTIVE_BORDER_COLOR = tcell.ColorOrange
	ICON_LOADING             = "\uea75 "
)

func CreateApp(prs []types.PR) *tview.Application {
	app := tview.NewApplication()

	// Get the current Git directory
	workspace, repoSlug, _ := util.GetRepoAndWorkspace()

	// UI components
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[::b]Bitbucket PR Viewer - %s - %s", workspace, repoSlug))

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	prDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select a PR to view details")

	activityDetails := tview.NewFlex().
		SetDirection(tview.FlexRow)

	diffDetails := tview.NewFlex().
		SetDirection(tview.FlexRow)

	diffStatDetails := tview.NewFlex().
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

	rightPanelHeader := tview.NewTextView().
		SetText("Here we will display PR title").
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	rightPanelGrid.AddItem(rightPanelHeader, 0, 0, 1, 1, 0, 0, false)
	rightPanelGrid.AddItem(prDetails, 1, 0, 1, 1, 0, 0, false)
	rightPanelGrid.AddItem(activityDetails, 2, 0, 1, 1, 0, 0, false)
	rightPanelGrid.AddItem(diffStatDetails, 1, 1, 1, 2, 0, 0, false)
	rightPanelGrid.AddItem(diffDetails, 2, 1, 1, 2, 0, 0, false)

	mainGrid.AddItem(rightPanelGrid, 1, 1, 1, 1, 0, 0, false)
	// Populate PR list
	populatePRList(prs, prList)

	prList.SetSelectedFunc(func(row, column int) {
		pr.UpdatePrDetails(prs, prDetails, row)

		if row >= 0 && row < len(prs) {
			selectedPR := prs[row]
			go func() {
				rightPanelHeader.SetText(FormatPRHeader(selectedPR))

				activityDetails.Clear()
				activityDetails.AddItem(tview.NewTextView().SetText(ICON_LOADING+"Fetching activities..."), 0, 1, true)

				diffDetails.Clear()
				diffDetails.AddItem(tview.NewTextView().SetText(ICON_LOADING+"Fetching diff..."), 0, 1, true)

				diffStatDetails.Clear()
				diffStatDetails.AddItem(tview.NewTextView().SetText(ICON_LOADING+"Fetching diff stats..."), 0, 1, true)

				diffData := fetchBitbucketDiff(selectedPR.ID)
				diffStatData := fetchBitbucketDiffstat(selectedPR.ID)
				prActivities := fetchBitbucketActivities(selectedPR.ID)

				app.QueueUpdateDraw(func() {
					activityDetails.Clear()
					activityDetails.AddItem(pr.CreateActivitiesView(prActivities), 0, 1, true)

					diffDetails.Clear()
					diffDetails.AddItem(pr.GenerateDiffView(diffData), 0, 1, true)

					diffStatDetails.Clear()
					diffStatDetails.AddItem(pr.GenerateDiffStatTree(diffStatData), 0, 1, true)
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
			activityDetails.Clear()
			rightPanelHeader.SetText(FormatPRHeader(initialPR))
			diffData := fetchBitbucketDiff(initialPR.ID)
			diffStatData := fetchBitbucketDiffstat(initialPR.ID)
			prActivities := fetchBitbucketActivities(initialPR.ID)

			app.QueueUpdateDraw(func() {
				activityDetails.Clear()

				activityDetails.AddItem(pr.CreateActivitiesView(prActivities), 0, 1, true)
				diffDetails.AddItem(pr.GenerateDiffView(diffData), 0, 1, true)
				diffStatDetails.AddItem(pr.GenerateDiffStatTree(diffStatData), 0, 1, true)
			})
		}()
	}

	app.SetRoot(mainGrid, true)
	// Capture the Tab key to switch focus between the views
	// Maintain a list of views in the desired focus order
	focusOrder := []tview.Primitive{prList, prDetails, activityDetails, diffStatDetails, diffDetails}
	currentFocusIndex := 0
	// Set initial borders
	util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
			app.SetFocus(focusOrder[currentFocusIndex])
			util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
		case tcell.KeyCtrlC:
			app.Stop()
		}
		return event
	})
	return app
}

// Function to populate the PR list
func populatePRList(prs []types.PR, prList *tview.Table) {
	for i, pr := range prs {
		titleCell := cellFormat(util.EllipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := util.CreateStateCell(pr.State)

		initialsCell := cellFormat(formatInitials(pr.Author.DisplayName), HIGH_CONTRAST_COLOR)

		sourceBranch := cellFormat(util.EllipsizeText(pr.Source.Branch.Name, 10), LOW_CONTRAST_COLOR)
		arrow := cellFormat("->", LOW_CONTRAST_COLOR)
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

// FormatPRHeader takes the PR details and returns a formatted string
func FormatPRHeader(pr types.PR) string {
	// Use fmt.Sprintf to format the header and apply tview's dynamic color syntax
	headerText := fmt.Sprintf(
		"%s\n\n"+
			"[yellow]%s[white] -> "+
			"[green]%s[white]",
		pr.Title,
		pr.Source.Branch.Name,
		pr.Destination.Branch.Name,
	)

	return headerText
}
