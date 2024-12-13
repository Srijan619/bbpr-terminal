// util/ui_utils.go
package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/state"
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
					SetBorderColor(activeBorderColor).
					SetBorderPadding(1, 1, 1, 1)
			} else {
				bordered.SetBorder(false)
			}
		}
	}
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

