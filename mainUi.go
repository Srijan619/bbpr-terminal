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
		SetBorderColor(tcell.ColorGrey).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Pull Requests(p)").SetBorderPadding(1, 1, 1, 1)

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	prListFlex.AddItem(prList, 0, 1, true)

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(prStatusFilterFlex, 0, 1, false).
		AddItem(prListFlex, 0, 18, true)

		// Description and Activity

	activityDetails := tview.NewFlex()
	activityDetails.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorGrey).
		SetTitle("Activities(a)").
		SetTitleAlign(tview.AlignLeft)

	// MIDDLE
	rightPanelHeader := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	rightPanelHeader.
		SetBorderColor(tcell.ColorGray).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle("Branch").
		SetTitleAlign(tview.AlignLeft)

	prDetails := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	prDetails.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorGrey).
		SetTitle("Description").
		SetTitleAlign(tview.AlignLeft)

	middleFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	middleFullFlex.AddItem(rightPanelHeader, 0, 1, false).
		AddItem(prDetails, 0, 2, false).
		AddItem(activityDetails, 0, 8, false)

		//RIGHT

	diffStatDetails := tview.NewFlex()
	diffStatDetails.SetTitle("Diff Tree(t)").
		SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetBorderPadding(1, 1, 1, 1).
		SetTitleAlign(tview.AlignLeft)

	diffDetails := tview.NewFlex()
	diffDetails.SetTitle("Diff Content(T)").
		SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetBorderPadding(1, 1, 1, 1).
		SetTitleAlign(tview.AlignLeft)
	rightFullFlex := tview.NewFlex()
	rightFullFlex.AddItem(diffStatDetails, 0, 1, false).
		AddItem(diffDetails, 0, 1, false)

	mainFlexWrapper := tview.NewFlex()

	mainFlexWrapper.AddItem(leftFullFlex, 60, 1, true).
		AddItem(middleFullFlex, 0, 1, false).
		AddItem(rightFullFlex, 0, 2, false)

	state.InitializeViews(app, mainFlexWrapper, prList, prDetails, activityDetails, diffDetails, diffStatDetails, prStatusFilterFlex, rightPanelHeader)
	pr.PopulatePRList(prList)

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
			case 't':
				state.GlobalState.App.SetRoot(state.GlobalState.DiffStatView, true)
			case 'T':
				state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
			case 'a':
				state.GlobalState.App.SetRoot(state.GlobalState.ActivityView, true)
			case 'q':
				state.GlobalState.App.SetRoot(state.GlobalState.MainFlexWrapper, true)
			case 'p':
				state.GlobalState.App.SetFocus(state.GlobalState.PrList)
			case 'm':
				state.SetPRStatusFilter("merged", !state.PRStatusFilter.Merged)
				//pr.UpdatePRList()
			case 'o':
				state.SetPRStatusFilter("open", !state.PRStatusFilter.Open)
			//	pr.UpdatePRList()
			case 'd':
				state.SetPRStatusFilter("declined", !state.PRStatusFilter.Declined)
				//pr.UpdatePRList() //TODO UI of filter is still not updating
			}
		}
		return event
	})
}
