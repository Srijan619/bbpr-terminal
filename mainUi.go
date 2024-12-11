package main

import (
	"github.com/rivo/tview"
	"simple-git-terminal/components/pr"
	"simple-git-terminal/util"
)

func CreateMainUi() *tview.Flex {
	prList := tview.NewFlex().
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Pull Requests")

	prStatusFilterFlex := pr.CreatePRStatusFilterView()

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(prStatusFilterFlex, 0, 1, false).
		AddItem(prList, 0, 18, false)

	middleFullFlex := util.CreateFlexComponent("Pull Request Details")
	rightFullFlex := util.CreateFlexComponent("Diff")

	mainFlexWrapper := tview.NewFlex()

	mainFlexWrapper.AddItem(leftFullFlex, 40, 1, true).
		AddItem(middleFullFlex, 0, 1, true).
		AddItem(rightFullFlex, 0, 2, true)

	return mainFlexWrapper
}
