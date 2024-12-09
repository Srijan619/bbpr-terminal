package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/pr"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

const (
	HIGH_CONTRAST_COLOR      = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR       = tcell.ColorYellow
	VIEW_ACTIVE_BORDER_COLOR = tcell.ColorOrange
	ICON_LOADING             = "\uea75 "
)

func CreateApp(prs []types.PR, workspace string, repoSlug string) *tview.Application {
	app := tview.NewApplication()

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

	activityDetails := tview.NewFlex()
	diffDetails := tview.NewFlex()
	diffStatDetails := tview.NewFlex()

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
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("Selected PR")

	rightPanelGrid.
		AddItem(rightPanelHeader, 0, 0, 1, 1, 0, 0, false).
		AddItem(prDetails, 1, 0, 1, 1, 0, 0, false).
		AddItem(activityDetails, 2, 0, 1, 1, 0, 0, false).
		AddItem(diffStatDetails, 1, 1, 1, 2, 0, 0, false).
		AddItem(diffDetails, 2, 1, 1, 2, 0, 0, false)

	mainGrid.AddItem(rightPanelGrid, 1, 1, 1, 1, 0, 0, false)
	// Populate PR list
	populatePRList(prs, prList)

	state.InitializeState(app, mainGrid, prList, prDetails, activityDetails, diffDetails, diffStatDetails)

	prList.SetSelectedFunc(func(row, column int) {
		pr.UpdatePrDetails(prs, prDetails, row)

		if row >= 0 && row < len(prs) {
			state.SetSelectedPR(&prs[row])
			go func() {
				rightPanelHeader.SetText(FormatPRHeader(*state.GlobalState.SelectedPR))

				util.UpdateActivityView(ICON_LOADING + "Fetching activities...")
				util.UpdateDiffStatView(ICON_LOADING + "Fetching diff stats...")
				util.UpdateDiffDetailsView("Select a file to see diff..")
				diffStatData := fetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)
				prActivities := fetchBitbucketActivities(state.GlobalState.SelectedPR.ID)

				app.QueueUpdateDraw(func() {
					util.UpdateActivityView(pr.CreateActivitiesView(prActivities))
					util.UpdateDiffStatView(pr.GenerateDiffStatTree(diffStatData))
				})
			}()
		}
	})

	//	Set initial PR details
	if len(prs) > 0 {
		prList.Select(0, 0)
		pr.UpdatePrDetails(prs, prDetails, 0)

		// Fetch initial activities dynamically
		state.SetSelectedPR(&prs[0])
		go func() {
			rightPanelHeader.SetText(FormatPRHeader(*state.GlobalState.SelectedPR))
			util.UpdateDiffDetailsView("Select a file to see diff..")
			diffStatData := fetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)
			prActivities := fetchBitbucketActivities(state.GlobalState.SelectedPR.ID)

			app.QueueUpdateDraw(func() {
				util.UpdateActivityView(pr.CreateActivitiesView(prActivities))
				util.UpdateDiffStatView(pr.GenerateDiffStatTree(diffStatData))
			})
		}()
	}

	// Key Bindings
	setupKeyBindings()

	app.SetRoot(mainGrid, true)

	return app
}

// Key bindings should be moved to somewher else later..
func setupKeyBindings() {
	// Capture the Tab key to switch focus between the views
	// Maintain a list of views in the desired focus order
	focusOrder := []tview.Primitive{state.GlobalState.PrList, state.GlobalState.PrDetails, state.GlobalState.ActivityView, state.GlobalState.DiffStatView, state.GlobalState.DiffDetails}
	currentFocusIndex := 0
	util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

	state.GlobalState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
			state.GlobalState.App.SetFocus(focusOrder[currentFocusIndex])
			util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
		case tcell.KeyCtrlC:
			state.GlobalState.App.Stop()
		case tcell.KeyRune:
			switch event.Rune() {
			case 'd':
				state.GlobalState.App.SetRoot(state.GlobalState.DiffStatView, true)
			case 'D':
				state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
			case 'a':
				state.GlobalState.App.SetRoot(state.GlobalState.ActivityView, true)
			case 'q':
				state.GlobalState.App.SetRoot(state.GlobalState.MainGrid, true)
			}
		}
		return event
	})
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
