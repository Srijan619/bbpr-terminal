// util/ui_utils.go
package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
)

const (
	HIGH_CONTRAST_COLOR = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR  = tcell.ColorYellow
	ICON_ACTIVE         = "\uf00c "
)

func GetStateColor(state string) tcell.Color {
	switch state {
	case "OPEN":
		return tcell.ColorLawnGreen
	case "MERGED":
		return tcell.ColorBlue
	case "DECLINED":
		return tcell.ColorRed
	default:
		return tcell.ColorYellow
	}
}

func GetFieldBasedColor(field string) tcell.Color {
	switch field {
	case "title":
		return tcell.ColorDarkCyan
	case "description":
		return tcell.ColorOrange
	default:
		return tcell.ColorWhite
	}
}

func cellFormat(text string, color tcell.Color) *tview.TableCell {
	return tview.NewTableCell(text).
		SetTextColor(color).
		SetAlign(tview.AlignLeft).
		SetSelectable(true)
}

// CreateStateCell creates a table cell with the appropriate color and alignment
func CreateStateCell(state string) *tview.TableCell {
	stateColor := GetStateColor(state)
	return tview.NewTableCell(state).
		SetTextColor(stateColor).
		SetAlign(tview.AlignLeft).
		SetSelectable(true)
}

func EllipsizeText(text string, max int) string {
	if len(text) > max {
		text = text[:max] + "..."
	}
	return text
}

func PopulatePRList(prList *tview.Table, prs []types.PR) {
	log.Printf("I am now population list view with PRs...%d", len(prs))
	for i, pr := range prs {
		titleCell := cellFormat(EllipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := CreateStateCell(pr.State)

		initialsCell := cellFormat(FormatInitials(pr.Author.DisplayName), HIGH_CONTRAST_COLOR)

		sourceBranch := cellFormat(EllipsizeText(pr.Source.Branch.Name, 10), LOW_CONTRAST_COLOR)
		arrow := cellFormat("->", LOW_CONTRAST_COLOR)
		destinationBranch := cellFormat(EllipsizeText(pr.Destination.Branch.Name, 10), LOW_CONTRAST_COLOR)

		activeCell := cellFormat(ICON_ACTIVE, tcell.ColorGreen)
		if state.GlobalState != nil && state.GlobalState.SelectedPR != nil && state.GlobalState.SelectedPR.ID == pr.ID {
			prList.SetCell(i, 0, activeCell)
		} else if i == 0 && (state.GlobalState == nil || state.GlobalState.SelectedPR == nil) {
			prList.SetCell(i, 0, activeCell)
		}
		prList.SetCell(i, 1, initialsCell)
		prList.SetCell(i, 2, stateCell)
		prList.SetCell(i, 3, titleCell)

		prList.SetCell(i, 4, sourceBranch)
		prList.SetCell(i, 5, arrow)
		prList.SetCell(i, 6, destinationBranch)
	}
}

// Helper method to update borders of views
func UpdateFocusBorders(focusOrder []tview.Primitive, currentFocusIndex int, activeBorderColor tcell.Color) {
	// for i, view := range focusOrder {
	// 	// Check if the view has border-related methods
	// 	// if bordered, ok := view.(interface {
	// 	// 	SetBorder(bool) *tview.Box
	// 	// 	SetBorderColor(tcell.Color) *tview.Box
	// 	// }); ok {
	// 	// 	if i == currentFocusIndex {
	// 	// 		// bordered.SetBorder(true).
	// 	// 		// 	SetBorderColor(activeBorderColor).
	// 	// 		// 	SetBorderPadding(1, 1, 1, 1)
	// 	// 	}
	// 	// }
	// }
}

func UpdateView(targetView *tview.Flex, content interface{}) {
	// Clear the target view before adding new content
	targetView.Clear()

	switch v := content.(type) {
	case string:
		textView := tview.NewTextView().
			SetText(v).
			SetDynamicColors(true).
			SetWrap(true)
		targetView.AddItem(textView, 0, 1, true)

	case tview.Primitive:
		targetView.AddItem(v, 0, 1, true)

	default:
		errorView := tview.NewTextView().
			SetText("[red]Unsupported content type[-]").
			SetDynamicColors(true).
			SetWrap(true)
		targetView.AddItem(errorView, 0, 1, true)
	}
}

func UpdateActivityView(activityContent interface{}) {
	UpdateView(state.GlobalState.ActivityView, activityContent)
}

func UpdateDiffDetailsView(diffContent interface{}) {
	UpdateView(state.GlobalState.DiffDetails, diffContent)
}

func UpdateDiffStatView(statContent interface{}) {
	UpdateView(state.GlobalState.DiffStatView, statContent)
}

func UpdatePRListView(prList []types.PR) {
	if state.GlobalState != nil && state.GlobalState.PrList != nil {
		PopulatePRList(state.GlobalState.PrList, prList)
	}
}
