// util/ui_utils.go
package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
)

const (
	HIGH_CONTRAST_COLOR = tcell.ColorCadetBlue
	LOW_CONTRAST_COLOR  = tcell.ColorYellow
	ICON_ACTIVE         = "\uf00c"
	ICON_SELECTED       = "\u25C8"
	ICON_DOWN_ARROW     = "\u2193"
	ICON_SIDE_ARROW     = "\u21AA"
	ICON_WARNING        = "\u2260"
	ICON_DECLINED       = "\u274C"
)

func GetPRStateColor(state string) tcell.Color {
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

func GetPRReviewStateIcon(state types.State) string {
	switch state {
	case types.StateApproved:
		return "[green]" + ICON_ACTIVE + "[-]"
	case types.StateDeclined:
		return "[red]" + ICON_DECLINED + "[-]"
	case types.StateRequestedChanges:
		return "[yellow]" + ICON_WARNING + "[-]"
	default:
		return ""
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
	stateColor := GetPRStateColor(state)
	return CreateTableCell(state, stateColor)
}

func EllipsizeText(text string, max int) string {
	if len(text) > max {
		text = text[:max] + "..."
	}
	return text
}

func PopulatePRList(prList *tview.Table, prs []types.PR) {
	// If there are no PRs, display a "No PRs" message
	if len(prs) == 0 {
		// Display a message in the first row
		noPRsCell := cellFormat("  No PRs available, try changing filters/search term", tcell.ColorWhite)
		prList.SetCell(0, 0, noPRsCell)
		return
	}

	for i, pr := range prs {
		titleCell := cellFormat(EllipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := CreateStateCell(pr.State)

		initialsCell := cellFormat(FormatInitials(pr.Author.DisplayName), HIGH_CONTRAST_COLOR)

		sourceBranch := cellFormat(EllipsizeText(pr.Source.Branch.Name, 18), tcell.ColorGrey)
		arrow := cellFormat(ICON_SIDE_ARROW, tcell.ColorDefault)
		destinationBranch := cellFormat(EllipsizeText(pr.Destination.Branch.Name, 18), tcell.ColorGrey)

		selectedCell := cellFormat(ICON_SELECTED, tcell.ColorOrange)
		if state.GlobalState != nil && state.GlobalState.SelectedPR != nil && state.GlobalState.SelectedPR.ID == pr.ID {
			prList.SetCell(i, 0, selectedCell)
		} else if i == 0 && (state.GlobalState == nil || state.GlobalState.SelectedPR == nil) {
			prList.SetCell(i, 0, selectedCell)
		}
		prList.SetCell(i, 1, initialsCell)
		prList.SetCell(i, 2, stateCell)

		prList.SetCell(i, 3, sourceBranch)
		prList.SetCell(i, 4, arrow)
		prList.SetCell(i, 5, destinationBranch)
		prList.SetCell(i, 6, titleCell)

	}

	prList.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkOrange))
}

// Helper method to update borders of views
func UpdateFocusBorders(focusOrder []tview.Primitive, currentFocusIndex int, activeBorderColor tcell.Color) {
	for i, view := range focusOrder {
		// Check if the view has border-related methods
		if bordered, ok := view.(interface {
			SetBorder(bool) *tview.Box
			SetBorderColor(tcell.Color) *tview.Box
		}); ok {
			if i == currentFocusIndex {
				bordered.SetBorder(true).
					SetBorderColor(activeBorderColor)
			} else {
				bordered.SetBorder(true).
					SetBorderColor(tcell.ColorGrey)
			}
		}
	}
}

func UpdateView(targetView interface{}, content interface{}) {
	if targetView != nil {
		// Check the type of the target view (either Flex or TextView)
		switch v := targetView.(type) {
		case *tview.Flex:
			// If it's a Flex view, clear it and update it
			v.Clear()

			// Handle content based on its type
			switch c := content.(type) {
			case string:
				// If the content is a string, display it in a TextView
				textView := CreateTextviewComponent("", false).SetText(c)
				v.AddItem(textView, 0, 1, true)
			case tview.Primitive:
				// If the content is a tview.Primitive, add it directly
				v.AddItem(c, 0, 1, true)
			default:
				// Handle unsupported content types
				errorView := CreateTextviewComponent("", false).SetText("[red]Unsupported content type[-]")
				v.AddItem(errorView, 0, 1, true)
			}

		case *tview.TextView:
			// If it's a TextView, update the text directly
			switch c := content.(type) {
			case string:
				// If content is a string, update the TextView
				v.SetText(c)
			case tview.Primitive:
				// Handle case if content is another Primitive (optional)
				log.Println("Unsupported Primitive content for TextView")
			default:
				// Handle unsupported content types
				v.SetText("[red]Unsupported content type[-]")
			}

		case *tview.Table:
			switch c := content.(type) {
			case string:
				tcell := CreateTableCell(c, tcell.ColorDefault)
				v.SetCell(0, 0, tcell)
			case tview.Primitive:
				// Handle case if content is another Primitive (optional)
				log.Println("Unsupported Primitive content for TextView")
			default:
				tcell := CreateTableCell("[red]Unsupported content type[-]", tcell.ColorDefault)
				v.SetCell(0, 0, tcell)
			}
		default:
			// If it's neither Flex nor TextView, print an error
			log.Println("[red]Unsupported target view type[-]")
		}
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

func UpdatePRListView() {
	if state.GlobalState != nil && state.GlobalState.PrList != nil && state.GlobalState.FilteredPRs != nil {
		state.GlobalState.PrList.Clear()
		PopulatePRList(state.GlobalState.PrList, *state.GlobalState.FilteredPRs)
		//	state.GlobalState.App.Draw()
	}
}

func UpdatePRListErrorView() {
	if state.GlobalState != nil && state.GlobalState.PrList != nil && state.GlobalState.FilteredPRs != nil {
		state.GlobalState.PrList.Clear()
		UpdateView(state.GlobalState.PrList, "[red] Error rendering PR list")
		state.GlobalState.App.Draw()
	}
}

func UpdatePRDetailView(content interface{}) {
	UpdateView(state.GlobalState.PrDetails, content)
}

func UpdatePRStatusFilterView(content interface{}) {
	if state.GlobalState != nil && state.GlobalState.PRStatusFilter != nil {
		UpdateView(state.GlobalState.PRStatusFilter, content)
	}
}

func UpdatePRListWithFilter(filter string, checked bool) {
	state.SetPRStatusFilter(filter, checked)
	ShowSpinnerFetchPRsByQueryAndUpdatePrList()
}

func UpdateFilteredPRs() {
	prs, pagination := bitbucket.FetchPRsByQuery(bitbucket.BuildQuery(""))
	state.SetFilteredPRs(&prs)
	state.SetPagination(&pagination)
}

func ShowSpinnerFetchPRsByQueryAndUpdatePrList() {
	if state.GlobalState != nil {
		state.GlobalState.PrList.Clear()
		ShowLoadingSpinner(state.GlobalState.PrList, func() (interface{}, error) {
			prs, pagination := bitbucket.FetchPRsByQuery(bitbucket.BuildQuery(state.SearchTerm))
			return struct {
				PRs        []types.PR
				Pagination types.Pagination
			}{prs, pagination}, nil
		}, func(result interface{}, err error) {
			if err != nil {
				UpdatePRListErrorView()
			} else {
				data := result.(struct {
					PRs        []types.PR
					Pagination types.Pagination
				})

				state.SetFilteredPRs(&data.PRs)
				state.SetPagination(&data.Pagination)
				UpdatePRListView()
				UpdatePaginationViewUI(state.Pagination.Page)
			}
		})
	}
}
