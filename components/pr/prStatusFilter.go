package pr

import (
	"github.com/rivo/tview"

	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	state.InitializePRStatusFilter(nil)

	openPr := util.CreateCheckBoxComponent("OPEN ").SetChecked(true)
	mergedPr := util.CreateCheckBoxComponent("MERGED ")
	declinedPr := util.CreateCheckBoxComponent("DECLINED ")
	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitle("Filter PRs by status").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}
