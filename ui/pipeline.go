package ui

import (
	"fmt"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func PopulatePPList(ppList *tview.Table, pps []types.PipelineResponse, frame int) {
	if len(pps) == 0 {
		noPPsCell := util.CellFormat("  No Pipelines available, try changing filters/search term", tcell.ColorWhite)
		ppList.SetCell(0, 0, noPPsCell)
		return
	}

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

		ppList.SetCell(i, 0, util.CellFormat(util.FormatInitials(pp.Creator.DisplayName), util.HIGH_CONTRAST_COLOR))      // Initial
		ppList.SetCell(i, 1, util.CellFormat(fmt.Sprintf("%s %d", util.ICON_BUILD, pp.BuildNumber), tcell.ColorDarkGray)) // Build #
		ppList.SetCell(i, 3, util.CellFormat(fmt.Sprintf(" %s", shortHash), tcell.ColorDarkGray))                        // Commit
		ppList.SetCell(i, 4, util.CellFormat(fmt.Sprintf(" %s", pp.Target.RefName), tcell.ColorDarkGray))                // Branch
		ppList.SetCell(i, 6, util.CellFormat(fmt.Sprintf("%s %s", statusIcon, status), statusColor))                      // Status
		ppList.SetCell(i, 9, util.CellFormat(durationStr, tcell.ColorDarkGray))                                           // Duration
		ppList.SetCell(i, 11, util.CellFormat(startStr, tcell.ColorDarkGray))                                             // Started

	}

	ppList.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkOrange))
}
