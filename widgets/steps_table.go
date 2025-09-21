package widgets

import (
	"fmt"
	"log"
	"simple-git-terminal/constants"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	widgets "simple-git-terminal/widgets/table"

	"github.com/gdamore/tcell/v2"
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
	stepTable.Clear()

	if len(steps) == 0 {
		stepTable.SetCell(0, 0, util.CellFormat(" No steps available", tcell.ColorGray))
		return
	}

	for i, step := range steps {
		status := step.State.Result.Name
		if status == "" {
			status = step.State.Name
		}

		color := util.GetColorForStatus(status)
		colorHex := util.HexColor(color)

		// Icon cell (col 0)
		var statusIcon string
		if status.NeedsTracking() {
			statusIcon = util.GetIconForStatusWithColorAnimated(status, frame)
		} else {
			statusIcon = util.GetIconForStatusWithColor(status)
		}

		if i == stepTable.SelectedRow {
			selectedCell := util.CellFormat(constants.ICON_SELECTED, tcell.ColorOrange)
			stepTable.SetCell(i, 0, selectedCell)
		} else {
			// Clear selection icon for other rows
			stepTable.SetCell(i, 0, util.CellFormat("", tcell.ColorDefault))
		}

		iconCell := util.CellFormat(fmt.Sprintf("[#%s:-]%s[-:-]", colorHex, statusIcon), color)
		stepTable.SetCell(i, 1, iconCell)

		// Name cell (col 1)
		nameCell := util.CellFormat(step.Name, tcell.ColorWhite)
		stepTable.SetCell(i, 2, nameCell)

	}
}

func (stepTable *StepsTable) UpdateSelectedRow(row int) {
	if row < 0 || row >= len(stepTable.steps) {
		return
	}

	// Redraw selection icon on all rows
	for i := range stepTable.steps {
		if i == row {
			stepTable.BaseTableView.UpdateSelectedRow(i)
		} else {
			stepTable.UpdateUnSelectedRow(i)
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

func (stepTable *StepsTable) PatchSteps(newSteps []types.StepDetail, frame int) {
	// Map current steps by UUID for quick lookup
	currentStepsMap := make(map[string]types.StepDetail)
	for _, step := range stepTable.steps {
		currentStepsMap[step.UUID] = step
	}

	for row, newStep := range newSteps {
		oldStep, found := currentStepsMap[newStep.UUID]

		if !found {
			// New step found — redraw whole table for simplicity
			log.Printf("[PatchSteps] New step detected (UUID: %s) at row %d, redrawing entire table", newStep.UUID, row)
			stepTable.SetSteps(newSteps, frame)
			return
		}

		// Determine old and new statuses to compare
		oldStatus := oldStep.State.Result.Name
		if oldStatus == "" {
			oldStatus = oldStep.State.Name
		}
		newStatus := newStep.State.Result.Name
		if newStatus == "" {
			newStatus = newStep.State.Name
		}

		log.Printf("[PatchSteps] Comparing step UUID: %s, row: %d", newStep.UUID, row)
		log.Printf("[PatchSteps] Old status: %s, New status: %s", oldStatus, newStatus)

		// Update icon if status changed or if status needs animation (e.g. in progress)
		if oldStatus != newStatus || newStatus.NeedsTracking() {
			color := util.GetColorForStatus(newStatus)
			colorHex := util.HexColor(color)

			var statusIcon string
			if newStatus.NeedsTracking() {
				statusIcon = util.GetIconForStatusWithColorAnimated(newStatus, frame)
			} else {
				statusIcon = util.GetIconForStatusWithColor(newStatus)
			}

			log.Printf("[PatchSteps] Updating icon for step UUID: %s at row %d", newStep.UUID, row)

			iconCell := util.CellFormat(fmt.Sprintf("[#%s:-]%s[-:-]", colorHex, statusIcon), color)
			stepTable.SetCell(row, 1, iconCell)
		} else {
			log.Printf("[PatchSteps] Icon for step UUID: %s unchanged", newStep.UUID)
		}

		// Update name cell if the step name changed
		if oldStep.Name != newStep.Name {
			log.Printf("[PatchSteps] Step name changed for UUID: %s, updating name cell at row %d", newStep.UUID, row)
			nameCell := util.CellFormat(newStep.Name, tcell.ColorWhite)
			stepTable.SetCell(row, 2, nameCell)
		} else {
			log.Printf("[PatchSteps] Name for step UUID: %s unchanged", newStep.UUID)
		}
	}

	// Save new steps state
	stepTable.steps = newSteps
	log.Printf("[PatchSteps] Completed patching steps, total steps: %d", len(newSteps))
}
