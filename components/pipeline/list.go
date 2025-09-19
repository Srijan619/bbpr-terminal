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
	// Validate row index
	if row < 0 || row >= len(pipelines) {
		log.Printf("Invalid row index: %d, pipelines count: %d", row, len(pipelines))
		return
	}

	if state.PipelineUIState == nil {
		log.Println("PipelineUIState is nil, cannot update UI")
		return
	}
	if state.PipelineUIState.PipelineDetails == nil {
		log.Println("PipelineDetails view is nil, cannot update UI")
		return
	}

	selectedPipeline := pipelines[row]
	log.Printf("Selected pipeline UUID: %s, Name: %s", selectedPipeline.UUID, selectedPipeline.Creator.DisplayName)

	util.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineDetails, func() (interface{}, error) {
		log.Println("Fetching pipeline steps...")
		steps := bitbucket.FetchPipelineSteps(selectedPipeline.UUID)
		if steps == nil {
			log.Println("Failed to fetch pipeline steps, nil returned")
			return nil, fmt.Errorf("failed to fetch pipeline steps")
		}
		log.Printf("Fetched %d steps", len(steps))
		return steps, nil
	}, func(result interface{}, err error) {
		log.Println("Stps...", result)
		steps, ok := result.([]types.StepDetail)
		if !ok {
			util.UpdateView(state.PipelineUIState.PipelineDetails, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		util.UpdateView(state.PipelineUIState.PipelineDetails, GeneratePPDetail(steps))
	})
}
