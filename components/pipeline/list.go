package pipeline

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/ui"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ICON_LOADING    = "\uea75 "
	ICON_SIDE_ARROW = "\u21AA "
)

func PopulatePipelineList(ppList *tview.Table) *tview.Table {
	pps, _ := bitbucket.FetchPipelinesByQuery(bitbucket.BuildQuery(""))
	// Populate PR list
	ui.PopulatePPList(ppList, pps)

	ppList.SetBackgroundColor(tcell.ColorDefault)

	// Populate pagination
	//
	// Add pagination below the PR list
	log.Printf("Current pagination state..%+v", state.Pagination)

	ppList.SetSelectedFunc(func(row, column int) {
		go func() {
			HandleOnPipelineSelect(pps, row)
		}()
	})
	return ppList
}

func HandleOnPipelineSelect(pipelines []types.PipelineResponse, row int) {
	// Ensure row index is valid
	if row < 0 || row >= len(pipelines) {
		log.Printf("Invalid row index: %d, pipeline count: %d", row, len(pipelines))
		return
	}

	if state.PipelineUIState == nil {
		log.Println("PipelineUIState is nil")
		return
	}

	selectedPipeline := pipelines[row]

	// Set header
	state.PipelineUIState.RightPanelHeader.SetTitle(fmt.Sprintf("Pipeline #%d - %s", selectedPipeline.BuildNumber, selectedPipeline.State.Result.Name))
	state.PipelineUIState.RightPanelHeader.SetText(fmt.Sprintf("[::b]Branch:[-:-] %s", selectedPipeline.Target.RefName))

	// Show loading spinner while fetching steps
	util.ShowLoadingSpinner(state.PipelineUIState.PipelineDetails, func() (interface{}, error) {
		// Fetch pipeline steps from bitbucket API
		steps := bitbucket.FetchPipelineSteps(selectedPipeline.UUID)
		if steps == nil {
			return nil, fmt.Errorf("Failed to fetch pipeline steps")
		}
		return &selectedPipeline, nil
	}, func(result interface{}, err error) {
		if err != nil {
			util.UpdateView(state.PipelineUIState.PipelineDetails, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		// Type assert result
		pipeline, ok := result.(*types.PipelineResponse)
		if !ok {
			util.UpdateView(state.PipelineUIState.PipelineDetails, "[red]Failed to cast pipeline details[-]")
			return
		}

		// Update view with generated pipeline detail
		util.UpdateView(state.PipelineUIState.PipelineDetails, GeneratePPDetail(pipeline))
	})
}
