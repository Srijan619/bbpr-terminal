package pr

import (
	"github.com/rivo/tview"

	"simple-git-terminal/components"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	openPr := components.CreateCheckBoxComponent("OPEN ")
	mergedPr := components.CreateCheckBoxComponent("MERGED ")
	declinedPr := components.CreateCheckBoxComponent("DECLINED ")
	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, true).
		SetTitle("Filter PRs by status").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	return wrapperFlex
}
