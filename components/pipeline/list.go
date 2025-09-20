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

func PopulatePipelineList(ppList *tview.Table) *tview.Table {
	ppList.SetBackgroundColor(tcell.ColorDefault)

	util.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineList, func() (interface{}, error) {
		pps, _ := bitbucket.FetchPipelinesByQuery(bitbucket.BuildQuery(""))
		if pps == nil {
			log.Println("Failed to fetch pipelines, nil returned")
			return nil, fmt.Errorf("failed to fetch pipelines")
		}
		return pps, nil
	}, func(result interface{}, err error) {
		pps, ok := result.([]types.PipelineResponse)
		if !ok {
			util.UpdateView(state.PipelineUIState.PipelineList, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}
		// Populate PR list
		ui.PopulatePPList(ppList, pps)

		ppList.SetSelectedFunc(func(row, column int) {
			go func() {
				HandleOnPipelineSelect(pps, row)
			}()
		})
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
	if state.PipelineUIState.PipelineSteps == nil {
		log.Println("PipelineSteps view is nil, cannot update UI")
		return
	}

	selectedPipeline := pipelines[row]
	log.Printf("Selected pipeline UUID: %s, Name: %s", selectedPipeline.UUID, selectedPipeline.Creator.DisplayName)

	util.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineSteps, func() (interface{}, error) {
		EmptyAllPipelineListDependentViews()

		steps := bitbucket.FetchPipelineSteps(selectedPipeline.UUID)
		if steps == nil {
			log.Println("Failed to fetch pipeline steps, nil returned")
			return nil, fmt.Errorf("failed to fetch pipeline steps")
		}
		return steps, nil
	}, func(result interface{}, err error) {
		steps, ok := result.([]types.StepDetail)
		if !ok {
			util.UpdateView(state.PipelineUIState.PipelineSteps, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		view := GenerateStepsView(steps, selectedPipeline)

		util.UpdateView(state.PipelineUIState.PipelineStepsDebugView, GeneratePPDebugInfo(selectedPipeline))
		util.UpdateView(state.PipelineUIState.PipelineSteps, view)

		HandleOnStepSelect(steps, selectedPipeline, 0) // Auto select first step and fetch the info
	})
}

func EmptyAllPipelineListDependentViews() {
	// Always start a fresh slate for rest of dependent views
	emptyView := tview.NewFlex()
	util.UpdateView(state.PipelineUIState.PipelineStepsDebugView, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineSteps, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineStep, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineStepCommandLogView, emptyView)
}
