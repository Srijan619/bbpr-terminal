package pr

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	if state.PRStatusFilter == nil {
		state.InitializePRStatusFilter(nil)
	}

	openPr := util.CreateCheckBoxComponent("Open (o) ", func(checked bool) {
		state.SetPRStatusFilter("open", checked)
		UpdatePRList()
	}).SetChecked(state.PRStatusFilter.Open)

	mergedPr := util.CreateCheckBoxComponent("Merged (m) ", func(checked bool) {
		state.SetPRStatusFilter("merged", checked)
		UpdatePRList()
	}).SetChecked(state.PRStatusFilter.Merged)

	declinedPr := util.CreateCheckBoxComponent("Declined (r) ", func(checked bool) {
		state.SetPRStatusFilter("declined", checked)
		UpdatePRList()
	}).SetChecked(state.PRStatusFilter.Declined)
	wrapperFlex.SetBackgroundColor(tcell.ColorDefault)
	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}

func UpdatePRList() {
	go func() {
		if state.GlobalState != nil {
			bitbucket.UpdateFilteredPRs()
			state.GlobalState.PrList.Clear()
			util.UpdatePRListView()
		}
	}()
}
