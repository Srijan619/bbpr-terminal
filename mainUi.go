package main

import (
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/components/pr"
	"simple-git-terminal/state"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
)

const (
	VIEW_ACTIVE_BORDER_COLOR = tcell.ColorOrange
)

func CreateMainApp() *tview.Application {
	app := tview.NewApplication()
	workspace, repoSlug, _ = util.GetRepoAndWorkspace()
	log.Printf("Loading workspace - %s and repo - %s ....", workspace, repoSlug)

	if (workspace == "") || (repoSlug == "") {
		log.Fatalf("Not a bitbucket Workspace")
	}
	state.SetWorkspaceRepo(workspace, repoSlug)
	//LEFT
	// PR Status Filter UI
	prStatusFilterFlex := pr.CreatePRStatusFilterView()

	// PR LIST UI
	prListFlex := tview.NewFlex()
	prListFlex.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Pull Requests")

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	pr.PopulatePRList(prList)
	prListFlex.AddItem(prList, 0, 1, true)

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(prStatusFilterFlex, 0, 1, false).
		AddItem(prListFlex, 0, 18, true)

		// Description and Activity

	prDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select a PR to view details")

	activityDetails := tview.NewFlex()

	// MIDDLE
	rightPanelHeader := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText("Selected PR")

	middleFullFlex := util.CreateFlexComponent("Pull Request Details")
	middleFullFlex.SetDirection(tview.FlexRow)

	middleFullFlex.AddItem(rightPanelHeader, 0, 1, false).
		AddItem(prDetails, 0, 2, false).
		AddItem(activityDetails, 0, 8, false)

		//RIGHT
	diffDetails := tview.NewFlex()
	diffStatDetails := tview.NewFlex()
	rightFullFlex := util.CreateFlexComponent("Diff")
	rightFullFlex.AddItem(diffStatDetails, 0, 1, false).
		AddItem(diffDetails, 0, 1, false)

	mainFlexWrapper := tview.NewFlex()

	mainFlexWrapper.AddItem(leftFullFlex, 60, 1, true).
		AddItem(middleFullFlex, 0, 1, false).
		AddItem(rightFullFlex, 0, 2, false)

	state.InitializeViews(app, mainFlexWrapper, prList, prDetails, activityDetails, diffDetails, diffStatDetails, rightPanelHeader)

	// Key Bindings
	setupKeyBindings()

	app.SetRoot(mainFlexWrapper, true).EnableMouse(true)

	return app
}

// Key bindings should be moved to somewher else later..
func setupKeyBindings() {
	// Capture the Tab key to switch focus between the views
	// Maintain a list of views in the desired focus order
	focusOrder := []tview.Primitive{state.GlobalState.PrList, state.GlobalState.PrDetails, state.GlobalState.ActivityView, state.GlobalState.DiffStatView, state.GlobalState.DiffDetails}
	currentFocusIndex := 0
	util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

	state.GlobalState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
			state.GlobalState.App.SetFocus(focusOrder[currentFocusIndex])
			util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
		case tcell.KeyCtrlC:
			state.GlobalState.App.Stop()
		case tcell.KeyRune:
			switch event.Rune() {
			case 'd':
				state.GlobalState.App.SetRoot(state.GlobalState.DiffStatView, true)
			case 'D':
				state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
			case 'a':
				state.GlobalState.App.SetRoot(state.GlobalState.ActivityView, true)
			case 'q':
				state.GlobalState.App.SetRoot(state.GlobalState.MainFlexWrapper, true)
			}
		}
		return event
	})
}
