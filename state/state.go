package state

import (
	"github.com/rivo/tview"

	"simple-git-terminal/types"
)

type State struct {
	App          *tview.Application
	MainGrid     *tview.Grid
	PrList       *tview.Table
	PrDetails    *tview.TextView
	ActivityView *tview.Flex
	DiffDetails  *tview.Flex
	DiffStatView *tview.Flex

	SelectedPR *types.PR
}

var GlobalState *State

func InitializeState(app *tview.Application, mainGrid *tview.Grid, prList *tview.Table, prDetails *tview.TextView, activityView, diffDetails, diffStatView *tview.Flex) {
	GlobalState = &State{
		App:          app,
		MainGrid:     mainGrid,
		PrList:       prList,
		PrDetails:    prDetails,
		ActivityView: activityView,
		DiffDetails:  diffDetails,
		DiffStatView: diffStatView,
	}
}

func SetSelectedPR(pr *types.PR) {
	GlobalState.SelectedPR = pr
}
