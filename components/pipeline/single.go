package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GenerateStepView(steps []types.StepDetail, selectedPipeline types.PipelineResponse) *tview.List {
	list := tview.NewList()

	for _, step := range steps {
		status := step.State.Result.Name
		if status == "" {
			status = step.State.Name
		}

		icon := util.GetIconForStatus(status)
		color := util.GetColorForStatus(status)
		colorHex := util.HexColor(color)

		title := fmt.Sprintf("[#%s]%s[-] %s", colorHex, icon, step.Name)

		list.AddItem(title, "", 0, nil)
	}

	list.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle(fmt.Sprintf(" Build %d - Steps ", selectedPipeline.BuildNumber)).
		SetTitleAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault)

	// Determine overall border color based on status of steps
	var borderColor tcell.Color = tcell.ColorDarkGray
	for _, step := range steps {
		status := step.State.Result.Name
		if status == "" {
			status = step.State.Name
		}

		if status.Failed() {
			borderColor = tcell.ColorRed
			break
		}
		if status.Successful() {
			borderColor = tcell.ColorGreen
		}
	}

	list.SetBorderColor(borderColor)

	return list
}
