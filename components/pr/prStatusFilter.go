package pr

import (
	"github.com/rivo/tview"

	"simple-git-terminal/components"
	"simple-git-terminal/state"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	state.InitializePRStatusFilter(nil)

	openPr := components.CreateCheckBoxComponent("OPEN ").SetChecked(true)
	mergedPr := components.CreateCheckBoxComponent("MERGED ")
	declinedPr := components.CreateCheckBoxComponent("DECLINED ")
	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitle("Filter PRs by status").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}
