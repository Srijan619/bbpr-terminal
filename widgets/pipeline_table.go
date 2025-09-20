package widgets

import (
	"fmt"
	"simple-git-terminal/constants"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PipelineTable struct {
	*tview.Table
	pipelines   []types.PipelineResponse
	selectedRow int
}

func NewPipelineTable() *PipelineTable {
	table := tview.NewTable()

	table.
		SetSelectable(true, false).
		SetFixed(1, 0).
		SetBackgroundColor(tcell.ColorDefault)

	return &PipelineTable{
		Table:       table,
		pipelines:   nil,
		selectedRow: -1,
	}
}

func (pt *PipelineTable) SetPipelines(pps []types.PipelineResponse, frame int) {
	pt.pipelines = pps

	// Clear and populate
	pt.Clear()

	for i, pp := range pps {
		shortHash := pp.Target.Commit.Hash
		if len(shortHash) > 7 {
			shortHash = shortHash[:7]
		}

		// If status result is not ready yet then it is ongoing...
		status := pp.State.Result.Name
		if status == "" {
			status = pp.State.Name
		}

		statusColor := util.GetColorForStatus(status)

		// Animated icon if in progress
		var statusIcon string
		if status.NeedsTracking() {
			statusIcon = util.GetIconForStatusWithColorAnimated(status, frame)
		} else {
			statusIcon = util.GetIconForStatusWithColor(status)
		}

		// Duration
		durationStr := fmt.Sprintf(" %ds", pp.Duration) // clock icon

		// Started Time - Human readable
		startTime, err := time.Parse(time.RFC3339Nano, pp.CreatedOn)
		var startStr string
		if err == nil {
			startStr = fmt.Sprintf(" %s", util.HumanizeTimeAgo(startTime)) // clock icon
		} else {
			startStr = " Unknown"
		}

		selectedCell := util.CellFormat(constants.ICON_SELECTED, tcell.ColorOrange)

		// no need to check what is selcted at this point, as this is very first time, select first row already
		if i == 0 {
			pt.SetCell(i, 0, selectedCell)
		}

		pt.SetCell(i, 1, util.CellFormat(util.FormatInitials(pp.Creator.DisplayName), constants.HIGH_CONTRAST_COLOR))      // Initial
		pt.SetCell(i, 2, util.CellFormat(fmt.Sprintf("%s %d", constants.ICON_BUILD, pp.BuildNumber), tcell.ColorDarkGray)) // Build #
		pt.SetCell(i, 4, util.CellFormat(fmt.Sprintf(" %s", shortHash), tcell.ColorDarkGray))                             // Commit
		pt.SetCell(i, 5, util.CellFormat(fmt.Sprintf(" %s", pp.Target.RefName), tcell.ColorDarkGray))                     // Branch
		pt.SetCell(i, 6, util.CellFormat(fmt.Sprintf("%s %s", statusIcon, status), statusColor))                           // Status
		pt.SetCell(i, 9, util.CellFormat(durationStr, tcell.ColorDarkGray))                                                // Duration
		pt.SetCell(i, 11, util.CellFormat(startStr, tcell.ColorDarkGray))                                                  // Started
	}
	pt.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkOrange))
}

func (pt *PipelineTable) UpdateSelectedRow(row int) {
	if row < 0 || row >= len(pt.pipelines) {
		return
	}

	pt.selectedRow = row

	// Redraw selection icon on all rows
	for i := 0; i < len(pt.pipelines); i++ {
		if i == row {
			pt.SetCell(i, 0, util.CellFormat(constants.ICON_SELECTED, tcell.ColorOrange))
		} else {
			pt.SetCell(i, 0, util.CellFormat("", tcell.ColorDefault))
		}
	}
}

func (pt *PipelineTable) GetSelectedPipeline() *types.PipelineResponse {
	if pt.selectedRow >= 0 && pt.selectedRow < len(pt.pipelines) {
		return &pt.pipelines[pt.selectedRow]
	}
	return nil
}

func (pt *PipelineTable) UpdateStatus(pipelineUUID string, status *tview.TableCell) {
	for i, pipeline := range pt.pipelines {
		if pipeline.UUID == pipelineUUID {
			pt.SetCell(i, 6, status)
		}
	}
}
