package pipeline

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/support"
	"simple-git-terminal/types"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GenerateStepsView(steps []types.StepDetail, selectedPipeline types.PipelineResponse, frame int) tview.Primitive {
	stepsTable := state.PipelineUIState.PipelineSteps

	stepsTable.
		SetBorders(false).
		SetBackgroundColor(tcell.ColorDefault)

	stepsTable.SetSteps(steps, frame)
	// Determine left bar color from most critical status
	barColor := util.GetColorForStatus(state.PipelineUIState.SelectedPipeline.State.Name) // Bar color is pipeline's state color

	listHeight := len(steps)
	if listHeight == 0 {
		listHeight = 1
	}

	// Left color bar
	leftBar := tview.NewBox().
		SetBorder(false).
		SetBackgroundColor(barColor)

	leftBarFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(leftBar, listHeight, 0, false)

	// Compose UI with left bar and step table
	layout := tview.NewFlex()

	layout.
		SetDirection(tview.FlexColumn).
		AddItem(leftBarFlex, 1, 0, false).
		AddItem(stepsTable, 0, 1, true).
		SetBackgroundColor(tcell.ColorDefault)

	stepsTable.SetSelectedFunc(func(row, column int) {
		go func() {
			HandleOnStepSelect(steps, selectedPipeline, row)
		}()
	})
	stepsTable.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkOrange))

	state.PipelineUIState.PipelineStepsTable = stepsTable
	return layout
}

func HandleOnStepSelect(steps []types.StepDetail, selectedPipeline types.PipelineResponse, row int) {
	// Validate row index
	if row < 0 || row >= len(steps) {
		log.Printf("Invalid row index: %d, steps count: %d", row, len(steps))
		return
	}

	if state.PipelineUIState == nil {
		log.Println("PipelineUIState is nil, cannot update UI")
		return
	}
	if state.PipelineUIState.PipelineStep == nil {
		log.Println("PipelineSteps view is nil, cannot update UI")
		return
	}

	selectedStep := steps[row]
	log.Printf("Selected step UUID: %s", selectedStep.UUID)
	state.PipelineUIState.PipelineSteps.UpdateSelectedRow(row)

	support.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineStep, func() (interface{}, error) {
		step := bitbucket.FetchPipelineStep(selectedPipeline.UUID, selectedStep.UUID)
		if step.UUID == "" {
			log.Println("Failed to fetch single step, empty UUID returned")
			return nil, fmt.Errorf("failed to fetch single step")
		}

		return step, nil
	}, func(result interface{}, err error) {
		step, ok := result.(types.StepDetail)
		if !ok {
			support.UpdateView(state.PipelineUIState.PipelineStep, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		view := GenerateStepView(step, selectedPipeline)

		support.UpdateView(state.PipelineUIState.PipelineStep, view)
	})
}
