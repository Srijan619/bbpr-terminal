package state

import (
	"log"
	"simple-git-terminal/types"
	"strings"

	"github.com/rivo/tview"
)

// InitializeViews initializes all view components except workspace and repo.

// PipelineState holds the global state for the pipeline UI.
type PipelineState struct {
	App                  *tview.Application
	MainFlexWrapper      *tview.Flex
	PipelineListFlex     *tview.Flex
	PipelineList         *tview.Table
	PipelineSteps        *tview.Flex
	PipelineStep         *tview.Flex
	RightPanelHeader     *tview.TextView
	CurrentView          tview.Primitive
	PipelineStatusFilter *tview.Flex
	PipelineSearchBar    *tview.InputField
	PaginationFlex       *tview.Flex

	SelectedPipeline  *types.PipelineResponse
	FilteredPipelines *[]types.PipelineResponse
}

// ✅ Unique name to avoid conflict with other state
var PipelineUIState *PipelineState

func InitializePipelineViews(
	app *tview.Application,
	mainFlexWrapper, pipelineListFlex *tview.Flex,
	pipelineList *tview.Table,
	pipelineSteps *tview.Flex,
	pipelineStep *tview.Flex,
	pipelineStatusFilter, paginationFlex *tview.Flex,
	pipelineSearchBar *tview.InputField,
) {
	PipelineUIState = &PipelineState{
		App:                  app,
		MainFlexWrapper:      mainFlexWrapper,
		PipelineListFlex:     pipelineListFlex,
		PipelineList:         pipelineList,
		PipelineSteps:        pipelineSteps,
		PipelineStep:         pipelineStep,
		PipelineStatusFilter: pipelineStatusFilter,
		PaginationFlex:       paginationFlex,
		PipelineSearchBar:    pipelineSearchBar,
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
