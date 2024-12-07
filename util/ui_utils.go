// util/ui_utils.go
package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetStateColor(state string) tcell.Color {
	switch state {
	case "OPEN":
		return tcell.ColorGreen
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
