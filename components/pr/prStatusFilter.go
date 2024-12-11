package pr

import (
	"github.com/rivo/tview"

	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	state.InitializePRStatusFilter(nil)

	openPr := util.CreateCheckBoxComponent("Open", func(checked bool) {
		state.SetPRStatusFilter("open", checked)
	})

	mergedPr := util.CreateCheckBoxComponent("Merged", func(checked bool) {
		state.SetPRStatusFilter("merged", checked)
	})

	declinedPr := util.CreateCheckBoxComponent("Declined", func(checked bool) {
		state.SetPRStatusFilter("declined", checked)
	})
	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitle("Filter PRs by status").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}
