package main

import (
	"github.com/rivo/tview"

	"simple-git-terminal/components"
	"simple-git-terminal/components/pr"
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

	middleFullFlex := components.CreateFlexComponent("Pull Request Details")
	rightFullFlex := components.CreateFlexComponent("Diff")

	mainFlexWrapper := tview.NewFlex()

	mainFlexWrapper.AddItem(leftFullFlex, 40, 1, true).
		AddItem(middleFullFlex, 0, 1, true).
		AddItem(rightFullFlex, 0, 2, true)

	return mainFlexWrapper
}

