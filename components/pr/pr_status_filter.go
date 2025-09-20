package pr

import (
	"simple-git-terminal/state"
	"simple-git-terminal/support"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreatePRStatusFilterView() *tview.Flex {
	wrapperFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	if state.PRStatusFilter == nil {
		state.InitializePRStatusFilter(nil)
	}

	checkboxes := []*tview.Checkbox{
		support.CreateCheckBoxComponent("Open (o) ", func(checked bool) {
			UpdatePRListWithFilter("open", checked)
		}).SetChecked(state.PRStatusFilter.Open),

		support.CreateCheckBoxComponent("Merged (m) ", func(checked bool) {
			UpdatePRListWithFilter("merged", checked)
		}).SetChecked(state.PRStatusFilter.Merged),

		support.CreateCheckBoxComponent("Declined (r) ", func(checked bool) {
			UpdatePRListWithFilter("declined", checked)
		}).SetChecked(state.PRStatusFilter.Declined),

		support.CreateCheckBoxComponent("I'm Author (I) ", func(checked bool) {
			UpdatePRListWithFilter("iamauthor", checked)
		}).SetChecked(state.PRStatusFilter.IAmAuthor),

		support.CreateCheckBoxComponent("I'm Reviewer (i) ", func(checked bool) {
			UpdatePRListWithFilter("iamreviewer", checked)
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

	// targetBranch := util.CreateDropDownComponent("Target Branch(t) => ", []string{"any", "First", "second"})
	// wrapperFlex.AddItem(targetBranch, 0, 1, false)
	wrapperFlex.SetBackgroundColor(tcell.ColorDefault).
		SetBorderPadding(0, 0, 1, 0).
		SetTitleAlign(tview.AlignLeft)

	return wrapperFlex
}
