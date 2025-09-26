package widgets

import (
	"fmt"

	"simple-git-terminal/constants"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	widgets "simple-git-terminal/widgets/table"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StepsCommandsTable struct {
	*widgets.BaseTableView
	step types.StepDetail
}

func NewStepsCommandsTable() *StepsCommandsTable {
	table := widgets.NewBaseTableView()

	table.SetTitle("Steps Info")

	return &StepsCommandsTable{
		BaseTableView: table,
	}
}

func (s *StepsCommandsTable) SetCell(row, col int, cell *tview.TableCell) {
	s.BaseTableView.SetCell(row, col, cell)
}

func (s *StepsCommandsTable) SetLoadingCell(cell *tview.TableCell) {
	s.BaseTableView.SetLoadingCell(cell)
}

func (stepTable *StepsCommandsTable) SetStepDetail(step types.StepDetail, frame int) {
	stepTable.ClearLoading()

	stepTable.step = step
	stepTable.Clear()

	// we store whole step, but this table only renders script commands
	scriptCommands := step.ScriptCommands
	if len(scriptCommands) == 0 {
		stepTable.SetCell(0, 0, util.CellFormat(" No script commands available", tcell.ColorGray))
		return
	}

	for i, cmd := range scriptCommands {
		cmdText := fmt.Sprintf("%s [::b]%s", constants.ICON_SIDE_ARROW, cmd.Name)
		stepTable.SetCell(i, 0, util.CellFormat(cmdText, tcell.ColorWhite))
	}

	if len(scriptCommands) == 0 {
		stepTable.SetCell(0, 0, util.CellFormat(" No script commands available", tcell.ColorGray))
	}
}

func (stepTable *StepsCommandsTable) UpdateSelectedRow(row int) {
	if row < 0 || row >= len(stepTable.step.ScriptCommands) {
		return
	}

	// Redraw selection icon on all rows
	for i := range stepTable.step.ScriptCommands {
		if i == row {
			stepTable.BaseTableView.UpdateSelectedRow(i)
		} else {
			stepTable.UpdateUnSelectedRow(i)
		}
	}
}

func (stepTable *StepsCommandsTable) UpdateStatus(command string, status *tview.TableCell) {
	for i, scriptCommand := range stepTable.step.ScriptCommands {
		if scriptCommand.Command == command {
			stepTable.SetCell(i, 6, status)
		}
	}
}
