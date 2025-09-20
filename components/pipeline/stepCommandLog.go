package pipeline

import (
	"fmt"
	"simple-git-terminal/util"

	"github.com/rivo/tview"
)

// GenerateStepCommandLogView renders the raw logs of a selected command step.
func GenerateStepCommandLogView(logText string, commandName string) tview.Primitive {
	logView := util.CreateTextviewComponent(fmt.Sprintf("   Logs: %s", commandName), false)
	logView.
		SetScrollable(true)

	if logText == "" {
		logView.SetText("[gray]No logs available for this command[-]")
	} else {
		logView.SetText(logText)
	}

	return logView
}
