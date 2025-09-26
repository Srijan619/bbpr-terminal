package state

import (
	"context"
	"log"
	"strings"

	"simple-git-terminal/types"
	"simple-git-terminal/widgets"

	"github.com/rivo/tview"
)

// InitializeViews initializes all view components except workspace and repo.

// PipelineState holds the global state for the pipeline UI.
type PipelineState struct {
	App                        *tview.Application
	MainFlexWrapper            *tview.Flex
	PipelineList               *widgets.PipelineTable
	PipelineStepsDebugView     *tview.Flex
	PipelineSteps              *widgets.StepsTable
	PipelineStep               *tview.Flex
	PipelineStepCommandsView   *widgets.StepsCommandsTable
	PipelineStepCommandLogView *tview.Flex
	RightPanelHeader           *tview.TextView
	CurrentView                tview.Primitive
	PipelineStatusFilter       *tview.Flex
	PipelineSearchBar          *tview.InputField
	PaginationFlex             *tview.Flex

	SelectedPipeline   *types.PipelineResponse
	FilteredPipelines  *[]types.PipelineResponse
	TrackingCancelFunc context.CancelFunc

	// dynamics
	PipelineStepsTable          *widgets.StepsTable
	PipelineScriptCommandsTable *tview.Table

	// central collection of views
	Views []tview.Primitive

	// for mocking
	IsNetworkMockMode bool
}

// ✅ Unique name to avoid conflict with other state
var PipelineUIState *PipelineState

func InitializePipelineViews(
	mainFlexWrapper *tview.Flex,
	pipelineList *widgets.PipelineTable,
	pipelineStepsDebugView *tview.Flex,
	pipelineSteps *widgets.StepsTable,
	pipelineStep *tview.Flex,
	pipelineStepCommandsView *widgets.StepsCommandsTable,
	pipelineStepCommandLogView *tview.Flex,
	pipelineStatusFilter, paginationFlex *tview.Flex,
	pipelineSearchBar *tview.InputField,
) {
	if PipelineUIState == nil {
		PipelineUIState = &PipelineState{}
	}

	// Only set fields if they aren't already set (allows partial init)
	if PipelineUIState.MainFlexWrapper == nil {
		PipelineUIState.MainFlexWrapper = mainFlexWrapper
	}
	if PipelineUIState.PipelineList == nil {
		PipelineUIState.PipelineList = pipelineList
	}
	if PipelineUIState.PipelineStepsDebugView == nil {
		PipelineUIState.PipelineStepsDebugView = pipelineStepsDebugView
	}
	if PipelineUIState.PipelineSteps == nil {
		PipelineUIState.PipelineSteps = pipelineSteps
	}
	if PipelineUIState.PipelineStep == nil {
		PipelineUIState.PipelineStep = pipelineStep
	}
	if PipelineUIState.PipelineStepCommandsView == nil {
		PipelineUIState.PipelineStepCommandsView = pipelineStepCommandsView
	}
	if PipelineUIState.PipelineStepCommandLogView == nil {
		PipelineUIState.PipelineStepCommandLogView = pipelineStepCommandLogView
	}
	if PipelineUIState.PipelineStatusFilter == nil {
		PipelineUIState.PipelineStatusFilter = pipelineStatusFilter
	}
	if PipelineUIState.PaginationFlex == nil {
		PipelineUIState.PaginationFlex = paginationFlex
	}
	if PipelineUIState.PipelineSearchBar == nil {
		PipelineUIState.PipelineSearchBar = pipelineSearchBar
	}

	// Update the views slice
	PipelineUIState.Views = []tview.Primitive{
		pipelineList,
		pipelineSteps,
		pipelineStepCommandsView,
		pipelineStatusFilter,
		pipelineSearchBar,
		pipelineStepCommandLogView,
		pipelineStep,
		pipelineStepsDebugView,
		paginationFlex,
		mainFlexWrapper,
	}
}

// PipelineStatusFilterType defines filters for pipeline states.
type PipelineStatusFilterType struct {
	Running  bool
	Success  bool
	Failed   bool
	Canceled bool
}

var PipelineStatusFilter *PipelineStatusFilterType

// InitializePipelineStatusFilter initializes the pipeline status filter with default values.
func InitializePipelineStatusFilter(filter *PipelineStatusFilterType) {
	if filter == nil {
		filter = &PipelineStatusFilterType{
			Running:  true,
			Success:  true,
			Failed:   true,
			Canceled: false,
		}
	}
	PipelineStatusFilter = filter
}

func InitPartialPipelineState(app *tview.Application, workspace, repo string) {
	if PipelineUIState == nil {
		PipelineUIState = &PipelineState{}
	}
	PipelineUIState.App = app
	Workspace = workspace
	Repo = repo
}

func SetSelectedPipeline(pipeline *types.PipelineResponse) {
	PipelineUIState.SelectedPipeline = pipeline
}

func updatePipelineListViewWithFreshFetch() {
}

// SetPipelineStatusFilter updates the pipeline status filter based on the provided key and value.
func SetPipelineStatusFilter(filterKey string, isChecked bool) {
	trimmedFilterKey := strings.ToLower(strings.TrimSpace(filterKey))
	switch trimmedFilterKey {
	case "running":
		PipelineStatusFilter.Running = isChecked
	case "success":
		PipelineStatusFilter.Success = isChecked
	case "failed":
		PipelineStatusFilter.Failed = isChecked
	case "canceled":
		PipelineStatusFilter.Canceled = isChecked
	case "all":
		PipelineStatusFilter.Running = isChecked
		PipelineStatusFilter.Success = isChecked
		PipelineStatusFilter.Failed = isChecked
		PipelineStatusFilter.Canceled = isChecked
	}
	log.Printf("Pipeline filter updated: %+v", PipelineStatusFilter)
}
