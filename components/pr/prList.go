package pr

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

const (
	HIGH_CONTRAST_COLOR = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR  = tcell.ColorYellow
	ICON_LOADING        = "\uea75 "
)

func PopulatePRList(prList *tview.Table) *tview.Table {
	prs := bitbucket.FetchBitbucketPRs()
	// Populate PR list
	populatePRList(prs)

	// Add a selection function that updates PR details when a PR is selected
	prList.SetSelectedFunc(func(row, column int) {
		HandleOnPrSelect(prs, row)
	})

	// Set initial PR details if available
	if len(prs) > 0 {
		prList.Select(0, 0)
		HandleOnPrSelect(prs, 0)
	}
	return prList
}

func HandleOnPrSelect(prs []types.PR, row int) {
	if row >= 0 && row < len(prs) {
		// Update PR details and set the selected PR
		UpdatePrDetails(prs, state.GlobalState.PrDetails, row)
		state.SetSelectedPR(&prs[row])

		// Update right panel and fetch diff stats and activities
		state.GlobalState.RightPanelHeader.SetText(formatPRHeader(*state.GlobalState.SelectedPR))

		// Fetch additional details in a separate goroutine to avoid blocking
		go func() {
			util.UpdateActivityView(ICON_LOADING + "Fetching activities...")
			util.UpdateDiffStatView(ICON_LOADING + "Fetching diff stats...")
			util.UpdateDiffDetailsView(" Hover over to a file for quick preview OR Select a file to see diff in full screen")

			// Fetch diff stat and activities for the selected PR
			diffStatData := bitbucket.FetchBitbucketDiffstat(state.GlobalState.SelectedPR.ID)
			prActivities := bitbucket.FetchBitbucketActivities(state.GlobalState.SelectedPR.ID)

			// Update views asynchronously
			state.GlobalState.App.QueueUpdateDraw(func() {
				util.UpdateActivityView(CreateActivitiesView(prActivities))
				util.UpdateDiffStatView(GenerateDiffStatTree(diffStatData))
			})
		}()
	}
}

// Function to populate the PR list
func populatePRList(prs []types.PR) {
	for i, pr := range prs {
		titleCell := cellFormat(util.EllipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := util.CreateStateCell(pr.State)

		initialsCell := cellFormat(util.FormatInitials(pr.Author.DisplayName), HIGH_CONTRAST_COLOR)

		sourceBranch := cellFormat(util.EllipsizeText(pr.Source.Branch.Name, 10), LOW_CONTRAST_COLOR)
		arrow := cellFormat("->", LOW_CONTRAST_COLOR)
		destinationBranch := cellFormat(util.EllipsizeText(pr.Destination.Branch.Name, 10), LOW_CONTRAST_COLOR)

		state.GlobalState.PrList.SetCell(i, 0, initialsCell)
		state.GlobalState.PrList.SetCell(i, 1, stateCell)
		state.GlobalState.PrList.SetCell(i, 2, titleCell)

		state.GlobalState.PrList.SetCell(i, 3, sourceBranch)
		state.GlobalState.PrList.SetCell(i, 4, arrow)
		state.GlobalState.PrList.SetCell(i, 5, destinationBranch)
	}
}

func cellFormat(text string, color tcell.Color) *tview.TableCell {
	return tview.NewTableCell(text).
		SetTextColor(color).
		SetAlign(tview.AlignLeft).
		SetSelectable(true)
}

// FormatPRHeader takes the PR details and returns a formatted string
func formatPRHeader(pr types.PR) string {
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
