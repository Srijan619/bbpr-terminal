package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GenerateStepView(step types.StepDetail, selectedPipeline types.PipelineResponse) tview.Primitive {
	table := tview.NewTable()

	table.
		SetBorders(false).
		SetSelectable(false, false).
		SetBackgroundColor(tcell.ColorDefault)

	// Determine status & color
	status := step.State.Name
	icon := util.GetIconForStatus(status)
	color := util.GetColorForStatus(status)
	colorHex := util.HexColor(color)

	// Metadata Rows
	row := 0
	table.SetCell(row, 0, util.CellFormat("[::b]Step UUID:[-]", tcell.ColorGray))
	table.SetCell(row, 1, util.CellFormat(step.UUID, tcell.ColorWhite))
	row++

	table.SetCell(row, 0, util.CellFormat("[::b]Status:[-]", tcell.ColorGray))
	table.SetCell(row, 1, util.CellFormat(fmt.Sprintf("[#%s]%s[-] %s", colorHex, icon, status), color))
	row++

	table.SetCell(row, 0, util.CellFormat("[::b]Started:[-]", tcell.ColorGray))
	table.SetCell(row, 1, util.CellFormat(util.FormatTime(step.StartedOn), tcell.ColorWhite))
	row++

	table.SetCell(row, 0, util.CellFormat("[::b]Completed:[-]", tcell.ColorGray))
	table.SetCell(row, 1, util.CellFormat(util.FormatTime(step.CompletedOn), tcell.ColorWhite))
	row++

	// Setup Commands
	if len(step.SetupCommands) > 0 {
		table.SetCell(row, 0, util.CellFormat("[::b]Setup Commands:[-]", tcell.ColorGray))
		row++
		for _, cmd := range step.SetupCommands {
			table.SetCell(row, 1, util.CellFormat(fmt.Sprintf("• [::b]%s:[-] %s", cmd.Name, cmd.Command), tcell.ColorWhite))
			row++
		}
	}

	// Script Commands
	if len(step.ScriptCommands) > 0 {
		table.SetCell(row, 0, util.CellFormat("[::b]Script Commands:[-]", tcell.ColorGray))
		row++
		for _, cmd := range step.ScriptCommands {
			table.SetCell(row, 1, util.CellFormat(fmt.Sprintf("• [::b]%s:[-] %s", cmd.Name, cmd.Command), tcell.ColorWhite))
			row++
		}
	}

	// Compose layout
	layout := tview.NewFlex()

	layout.
		SetDirection(tview.FlexColumn).
		AddItem(table, 0, 1, true).
		SetBackgroundColor(tcell.ColorDefault)

	return layout
}
