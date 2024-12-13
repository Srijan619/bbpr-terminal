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
		updatePRList()
	}).SetChecked(state.PRStatusFilter.Open)

	mergedPr := util.CreateCheckBoxComponent("Merged", func(checked bool) {
		state.SetPRStatusFilter("merged", checked)
		updatePRList()
	}).SetChecked(state.PRStatusFilter.Merged)

	declinedPr := util.CreateCheckBoxComponent("Declined", func(checked bool) {
		state.SetPRStatusFilter("declined", checked)
		updatePRList()
	}).SetChecked(state.PRStatusFilter.Declined)

	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitle("Filter PRs by status").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}

func updatePRList() {
	go func() {
		if state.GlobalState != nil {
			filteredPRs := GetFilteredPRs()
			state.GlobalState.PrList.Clear()
			state.GlobalState.App.QueueUpdateDraw(func() {
				util.UpdatePRListView(filteredPRs)
			})
		}
	}()
}
