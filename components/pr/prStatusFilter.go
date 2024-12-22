package pr

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	if state.PRStatusFilter == nil {
		state.InitializePRStatusFilter(nil)
	}

	checkboxes := []*tview.Checkbox{
		util.CreateCheckBoxComponent("Open (o) ", func(checked bool) {
			util.UpdatePRListWithFilter("open", checked)
		}).SetChecked(state.PRStatusFilter.Open),

		util.CreateCheckBoxComponent("Merged (m) ", func(checked bool) {
			util.UpdatePRListWithFilter("merged", checked)
		}).SetChecked(state.PRStatusFilter.Merged),

		util.CreateCheckBoxComponent("Declined (r) ", func(checked bool) {
			util.UpdatePRListWithFilter("declined", checked)
		}).SetChecked(state.PRStatusFilter.Declined),

		util.CreateCheckBoxComponent("I'm Author (I) ", func(checked bool) {
			util.UpdatePRListWithFilter("iamauthor", checked)
		}).SetChecked(state.PRStatusFilter.IAmAuthor),

		util.CreateCheckBoxComponent("I'm Reviewer (i) ", func(checked bool) {
			util.UpdatePRListWithFilter("iamreviewer", checked)
		}).SetChecked(state.PRStatusFilter.IAmReviewer),
	}

	rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	wrapperFlex.AddItem(rowFlex, 0, 1, false)

	maxWidth := 50
	currentWidth := 0

	for _, checkbox := range checkboxes {
		itemWidth := len(checkbox.GetLabel()) + 4
		if currentWidth+itemWidth > maxWidth {
			rowFlex = tview.NewFlex().SetDirection(tview.FlexColumn)
			wrapperFlex.AddItem(rowFlex, 0, 1, false)
			currentWidth = 0
		}

		rowFlex.AddItem(checkbox, itemWidth, 1, false)
		currentWidth += itemWidth

	}

	wrapperFlex.SetBackgroundColor(tcell.ColorDefault).
		SetBorderPadding(0, 0, 1, 0).
		SetTitleAlign(tview.AlignLeft)

	return wrapperFlex
}
