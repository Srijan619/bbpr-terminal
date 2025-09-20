package ui

import (
	"simple-git-terminal/constants"
	"simple-git-terminal/types"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func PopulatePRList(prList *tview.Table, prs []types.PR) {
	// If there are no PRs, display a "No PRs" message
	if len(prs) == 0 {
		// Display a message in the first row
		noPRsCell := util.CellFormat("  No PRs available, try changing filters/search term", tcell.ColorWhite)
		prList.SetCell(0, 0, noPRsCell)
		return
	}

	for i, pr := range prs {
		titleCell := util.CellFormat(util.EllipsizeText(pr.Title, 18), tcell.ColorWhite)
		stateCell := util.CreateStateCell(pr.State)

		initialsCell := util.CellFormat(util.FormatInitials(pr.Author.DisplayName), constants.HIGH_CONTRAST_COLOR)

		sourceBranch := util.CellFormat(util.EllipsizeText(pr.Source.Branch.Name, 18), tcell.ColorGrey)
		arrow := util.CellFormat(constants.ICON_SIDE_ARROW, tcell.ColorDefault)
		destinationBranch := util.CellFormat(util.EllipsizeText(pr.Destination.Branch.Name, 18), tcell.ColorGrey)

		selectedCell := util.CellFormat(constants.ICON_SELECTED, tcell.ColorOrange)

		// no need to check what is selcted at this point, as this is very first time, select first row already
		if i == 0 {
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
