package support

import (
	"simple-git-terminal/constants"
	"simple-git-terminal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

func GetPRReviewStateIcon(state types.ApprovedState) string {
	switch state {
	case types.StateApproved:
		return "[green]" + constants.ICON_ACTIVE + "[-]"
	case types.StateDeclined:
		return "[red]" + constants.ICON_DECLINED + "[-]"
	case types.StateRequestedChanges:
		return "[yellow]" + constants.ICON_WARNING + "[-]"
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

func CellFormat(text string, color tcell.Color) *tview.TableCell {
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
