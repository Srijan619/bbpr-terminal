package state

import (
	"github.com/rivo/tview"
	"simple-git-terminal/types"
)

type State struct {
	App              *tview.Application
	MainGrid         *tview.Grid
	PrList           *tview.Table
	PrDetails        *tview.TextView
	ActivityView     *tview.Flex
	DiffDetails      *tview.Flex
	DiffStatView     *tview.Flex
	RightPanelHeader *tview.TextView

	SelectedPR *types.PR
}

var GlobalState *State
var Workspace, Repo string

// InitializeViews initializes all view components except workspace and repo.
func InitializeViews(app *tview.Application, mainGrid *tview.Grid, prList *tview.Table, prDetails *tview.TextView, activityView, diffDetails, diffStatView *tview.Flex, rightPanelHeader *tview.TextView) {
	GlobalState = &State{
		App:              app,
		MainGrid:         mainGrid,
		PrList:           prList,
		PrDetails:        prDetails,
		ActivityView:     activityView,
		DiffDetails:      diffDetails,
		DiffStatView:     diffStatView,
		RightPanelHeader: rightPanelHeader,
	}
}

// SetWorkspaceRepo sets the workspace and repo separately from the main state
func SetWorkspaceRepo(workspace, repo string) {
	Workspace = workspace
	Repo = repo
}

// SetSelectedPR sets the selected PR in the global state.
func SetSelectedPR(pr *types.PR) {
	GlobalState.SelectedPR = pr
}
