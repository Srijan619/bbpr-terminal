package widgets

import (
	"simple-git-terminal/types"
	widgets "simple-git-terminal/widgets/table"

	"github.com/rivo/tview"
)

type StepsTable struct {
	*widgets.BaseTableView
	steps []types.StepDetail
}

func NewStepsTable() *StepsTable {
	table := widgets.NewBaseTableView()

	table.SetTitle("Steps Info").
		SetTitleAlign(tview.AlignLeft)

	return &StepsTable{
		BaseTableView: table,
		steps:         nil,
	}
}

func (stepTable *StepsTable) SetSteps(steps []types.StepDetail, frame int) {
	stepTable.steps = steps

	// Clear and populate
	stepTable.Clear()
}

func (stepTable *StepsTable) UpdateSelectedRow(row int) {
	if row < 0 || row >= len(stepTable.steps) {
		return
	}

	// Redraw selection icon on all rows
	for i := range len(stepTable.steps) {
		if i == row {
			stepTable.BaseTableView.UpdateSelectedRow(row)
		} else {
			stepTable.BaseTableView.UpdateUnSelectedRow(row)
		}
	}
}

func (stepTable *StepsTable) GetSelectedStep() *types.StepDetail {
	if stepTable.SelectedRow >= 0 && stepTable.SelectedRow < len(stepTable.steps) {
		return &stepTable.steps[stepTable.SelectedRow]
	}
	return nil
}

func (stepTable *StepsTable) UpdateStatus(stepUUID string, status *tview.TableCell) {
	for i, step := range stepTable.steps {
		if step.UUID == stepUUID {
			stepTable.SetCell(i, 6, status)
		}
	}
}
