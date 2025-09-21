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
	stepTable := tview.NewTable()

	stepTable.
		SetBorders(false).
		SetBackgroundColor(tcell.ColorDefault)

		// Determine left bar color from most critical status
	barColor := tcell.ColorDarkGray

	if len(steps) == 0 {
		stepTable.SetCell(0, 0, util.CellFormat(" No steps available", tcell.ColorGray))
	} else {
		for i, step := range steps {
			status := step.State.Result.Name
			if status == "" {
				status = step.State.Name
			}

			color := util.GetColorForStatus(status)
			colorHex := util.HexColor(color)

			// Animated icon if in progress
			var statusIcon string
			if status.NeedsTracking() {
				statusIcon = util.GetIconForStatusWithColorAnimated(status, frame)
			} else {
				statusIcon = util.GetIconForStatusWithColor(status)
			}

			barColor = color
			text := fmt.Sprintf(" [#%s:-]%s[-:-] %s", colorHex, statusIcon, step.Name)
			cell := util.CellFormat(text, color)
			stepTable.SetCell(i, 0, cell)
		}
	}

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
		AddItem(stepTable, 0, 1, true).
		SetBackgroundColor(tcell.ColorDefault)

	stepTable.SetSelectedFunc(func(row, column int) {
		go func() {
			HandleOnStepSelect(steps, selectedPipeline, row)
		}()
	})
	stepTable.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkOrange))

	state.PipelineUIState.PipelineStepTable = stepTable
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
