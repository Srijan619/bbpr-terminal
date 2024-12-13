package main

import (
	"fmt"
	"log"

	"github.com/rivo/tview"

	"simple-git-terminal/components/pr"
	"simple-git-terminal/state"
	"simple-git-terminal/util"
)

func CreateApp() *tview.Application {
	workspace, repoSlug, _ = util.GetRepoAndWorkspace()
	log.Printf("Loading workspace - %s and repo - %s ....", workspace, repoSlug)

	if (workspace == "") || (repoSlug == "") {
		log.Fatalf("Not a bitbucket Workspace")
	}
	state.SetWorkspaceRepo(workspace, repoSlug)
	app := tview.NewApplication()

	// UI components
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[::b]Bitbucket PR Viewer - %s - %s", workspace, repoSlug))

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	prDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select a PR to view details")

	activityDetails := tview.NewFlex()
	diffDetails := tview.NewFlex()
	diffStatDetails := tview.NewFlex()

	// Grid layout
	mainGrid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(60, 0).
		SetBorders(true)

	mainGrid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	mainGrid.AddItem(prList, 1, 0, 1, 1, 0, 0, true)
	rightPanelGrid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(60, 0).
		SetBorders(true)

	rightPanelHeader := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText("Selected PR")

	rightPanelGrid.
		AddItem(rightPanelHeader, 0, 0, 1, 2, 0, 0, false).
		AddItem(prDetails, 1, 0, 1, 1, 0, 0, false).
		AddItem(activityDetails, 2, 0, 1, 1, 0, 0, false).
		AddItem(diffStatDetails, 1, 1, 1, 2, 0, 0, false).
		AddItem(diffDetails, 2, 1, 1, 2, 0, 0, false)

	mainGrid.AddItem(rightPanelGrid, 1, 1, 1, 1, 0, 0, false)

	//TODO state.InitializeViews(app, mainGrid, prList, prDetails, activityDetails, diffDetails, diffStatDetails, rightPanelHeader)
	pr.PopulatePRList(prList)

	// Key Bindings
	setupKeyBindings()

	app.SetRoot(mainGrid, true)
	state.SetCurrentView(mainGrid)

	return app
}

// Key bindings should be moved to somewher else later..
// func setupKeyBindings() {
// 	// Capture the Tab key to switch focus between the views
// 	// Maintain a list of views in the desired focus order
// 	focusOrder := []tview.Primitive{state.GlobalState.PrList, state.GlobalState.PrDetails, state.GlobalState.ActivityView, state.GlobalState.DiffStatView, state.GlobalState.DiffDetails}
// 	currentFocusIndex := 0
// 	util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
//
// 	state.GlobalState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyTAB:
// 			currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
// 			state.GlobalState.App.SetFocus(focusOrder[currentFocusIndex])
// 			util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
// 		case tcell.KeyCtrlC:
// 			state.GlobalState.App.Stop()
// 		case tcell.KeyRune:
// 			switch event.Rune() {
// 			case 'd':
// 				state.SetCurrentView(state.GlobalState.DiffStatView)
// 				state.GlobalState.App.SetRoot(state.GlobalState.DiffStatView, true)
// 			case 'D':
// 				state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
// 			case 'a':
// 				state.GlobalState.App.SetRoot(state.GlobalState.ActivityView, true)
// 			case 'q':
// 				state.GlobalState.App.SetRoot(state.GlobalState.MainGrid, true)
// 				util.UpdateFocusBorders(focusOrder, 0, VIEW_ACTIVE_BORDER_COLOR)
// 			}
// 		}
// 		return event
// 	})
// }
