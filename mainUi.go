package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/components/pr"
	"simple-git-terminal/state"
	"simple-git-terminal/util"
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
	util.InitMdRenderer() // Markdown renderer takes time, so init it beforehand

	//LEFT

	// PR Status Filter UI
	prStatusFilterFlex := util.CreateFlexComponent("Filters")
	prStatusFilterFlex.AddItem(pr.CreatePRStatusFilterView(), 0, 1, false)

	// PR LIST UI
	prListFlex := util.CreateFlexComponent("Pull Requests  [green]p").
		SetDirection(tview.FlexRow)

	prList := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)

	prList.SetBackgroundColor(tcell.ColorDefault)

	prListSearchBar := util.CreateTextAreaComponent("🔎 Search PR (s)", "search filter.....")

	prListFlex.
		AddItem(prList, 0, 1, true)

	leftFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	leftFullFlex.
		AddItem(prStatusFilterFlex, 0, 1, false).
		AddItem(prListSearchBar, 0, 2, false).
		AddItem(prListFlex, 0, 18, true)

		// Description and Activity

	activityDetails := util.CreateFlexComponent("Activities [green]a|A")

	// MIDDLE
	rightPanelHeader := util.CreateTextviewComponent("", true)
	prDetails := util.CreateTextviewComponent("Description [green]d|D", true)

	middleFullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	middleFullFlex.SetBackgroundColor(tcell.ColorDefault)

	middleFullFlex.AddItem(rightPanelHeader, 0, 1, false).
		AddItem(prDetails, 0, 2, false).
		AddItem(activityDetails, 0, 10, false)

		//RIGHT

	diffStatDetails := util.CreateFlexComponent("Diff Tree [green]t|T")
	diffDetails := util.CreateFlexComponent("Diff Content [green]c|C")

	rightFullFlex := tview.NewFlex()

	rightFullFlex.SetBackgroundColor(tcell.ColorDefault)
	rightFullFlex.AddItem(diffStatDetails, 0, 1, false).
		AddItem(diffDetails, 0, 1, false)

	mainFlexWrapper := tview.NewFlex()
	mainFlexWrapper.SetBackgroundColor(tcell.ColorDefault)
	mainFlexWrapper.AddItem(leftFullFlex, 0, 1, true).
		AddItem(middleFullFlex, 0, 1, false).
		AddItem(rightFullFlex, 0, 2, false)

	state.InitializeViews(app, mainFlexWrapper, prListFlex, prList, prDetails, activityDetails, diffDetails, diffStatDetails, prStatusFilterFlex, rightPanelHeader, prListSearchBar)
	pr.PopulatePRList(prList)

	// Key Bindings
	setupKeyBindings()

	app.SetRoot(mainFlexWrapper, true).EnableMouse(true)

	return app
}

func setupKeyBindings() {
	// Define focus order
	focusOrder := []tview.Primitive{
		state.GlobalState.PrListFlex, state.GlobalState.PrDetails, state.GlobalState.ActivityView,
		state.GlobalState.DiffStatView, state.GlobalState.DiffDetails, state.GlobalState.PrListSearchBar,
	}
	currentFocusIndex := 0
	util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

	state.GlobalState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If in search mode, only allow Esc or Enter keys
		if state.IsSearchMode {
			switch event.Key() {
			case tcell.KeyEsc:
				state.SetIsSearchMode(false)
				state.GlobalState.App.SetFocus(state.GlobalState.PrList) // Focus back to PrList or another view
			case tcell.KeyEnter:
				// Perform search or exit search mode
				state.SetIsSearchMode(false)
				state.GlobalState.App.SetFocus(state.GlobalState.PrList) // Focus back to PrList or another view
			default:
				return event // Ignore other keys in search mode
			}
			util.UpdateFocusBorders(focusOrder, 0, VIEW_ACTIVE_BORDER_COLOR)

		} else {
			// Handle keybindings when not in search mode
			switch event.Key() {
			case tcell.KeyTAB:
				// Cycle focus between views
				currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
				state.GlobalState.App.SetFocus(focusOrder[currentFocusIndex])

			case tcell.KeyCtrlC:
				state.GlobalState.App.Stop()

			case tcell.KeyRune:
				switch event.Rune() {
				case 's':
					// Search mode
					state.SetIsSearchMode(true)
					currentFocusIndex = 5
					state.GlobalState.App.SetFocus(state.GlobalState.PrListSearchBar)

				case 't', 'T':
					// Focus on DiffStatView (T or t)
					currentFocusIndex = 3
					state.GlobalState.App.SetFocus(state.GlobalState.DiffStatView)

				case 'c', 'C':
					// Focus on DiffDetails (C or c)
					currentFocusIndex = 4
					state.GlobalState.App.SetFocus(state.GlobalState.DiffDetails)

				case 'a', 'A':
					// Focus on ActivityView (A or a)
					currentFocusIndex = 2
					state.GlobalState.App.SetFocus(state.GlobalState.ActivityView)

				case 'p':
					// Focus on PR List
					currentFocusIndex = 0
					state.GlobalState.App.SetFocus(state.GlobalState.PrList)

				case 'd', 'D':
					// Focus on PR Details (D or d)
					currentFocusIndex = 1
					state.GlobalState.App.SetFocus(state.GlobalState.PrDetails)

				case 'q':
					// Quit application
					state.GlobalState.App.SetRoot(state.GlobalState.MainFlexWrapper, true)

				case 'm', 'o', 'r':
					// Toggle PR filters
					switch event.Rune() {
					case 'm':
						state.SetPRStatusFilter("merged", !state.PRStatusFilter.Merged)
					case 'o':
						state.SetPRStatusFilter("open", !state.PRStatusFilter.Open)
					case 'r':
						state.SetPRStatusFilter("declined", !state.PRStatusFilter.Declined)
					}
					updateFilter()
				}
			}
			// Update focus borders after focus change
			util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
		}

		return event
	})
}

func updateFilter() {
	view := pr.CreatePRStatusFilterView()
	util.UpdatePRStatusFilterView(view)
}
