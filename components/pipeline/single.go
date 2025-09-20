package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GenerateStepView(steps []types.StepDetail, selectedPipeline types.PipelineResponse) tview.Primitive {
	stepTable := tview.NewTable()

	stepTable.
		SetBorders(false).
		SetSelectable(false, false).
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

			icon := util.GetIconForStatus(status)
			color := util.GetColorForStatus(status)
			colorHex := util.HexColor(color)

			barColor = color
			text := fmt.Sprintf("[#%s:-]%s[-:-] %s", colorHex, icon, step.Name)
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

	return layout
}
