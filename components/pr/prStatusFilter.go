package pr

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreatePRStatusFilterView() *tview.Flex {
	log.Printf("I am creating..%v", state.PRStatusFilter)
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	if state.PRStatusFilter == nil {
		log.Printf("Sont be the")
		state.InitializePRStatusFilter(nil)
	}

	openPr := util.CreateCheckBoxComponent("Open ", func(checked bool) {
		state.SetPRStatusFilter("open", checked)
		UpdatePRList()
	}).SetChecked(state.PRStatusFilter.Open)

	mergedPr := util.CreateCheckBoxComponent("Merged ", func(checked bool) {
		state.SetPRStatusFilter("merged", checked)
		UpdatePRList()
	}).SetChecked(state.PRStatusFilter.Merged)

	declinedPr := util.CreateCheckBoxComponent("Declined ", func(checked bool) {
		state.SetPRStatusFilter("declined", checked)
		UpdatePRList()
	}).SetChecked(state.PRStatusFilter.Declined)

	wrapperFlex.AddItem(openPr, 0, 1, false).
		AddItem(mergedPr, 0, 1, false).
		AddItem(declinedPr, 0, 1, false).
		SetTitle("Filter PRs by status").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tcell.ColorGrey).
		SetBorderPadding(0, 0, 1, 0)

	return wrapperFlex
}

func UpdatePRList() {
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
