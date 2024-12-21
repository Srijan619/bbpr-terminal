package pr

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	if state.PRStatusFilter == nil {
		state.InitializePRStatusFilter(nil)
	}

	openPr := util.CreateCheckBoxComponent("Open (o) ", func(checked bool) {
		util.UpdatePRListWithFilter("open", checked)
	}).SetChecked(state.PRStatusFilter.Open)

	mergedPr := util.CreateCheckBoxComponent("Merged (m) ", func(checked bool) {
		util.UpdatePRListWithFilter("merged", checked)
	}).SetChecked(state.PRStatusFilter.Merged)

	declinedPr := util.CreateCheckBoxComponent("Declined (r) ", func(checked bool) {
		util.UpdatePRListWithFilter("declined", checked)
	}).SetChecked(state.PRStatusFilter.Declined)
	wrapperFlex.SetBackgroundColor(tcell.ColorDefault)
	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}
