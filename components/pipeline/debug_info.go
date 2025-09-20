package pipeline

import (
	"fmt"
	"simple-git-terminal/constants"
	"simple-git-terminal/support"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GeneratePPDebugInfo(pipeline types.PipelineResponse) *tview.TextView {
	textView := support.CreateTextviewComponent("Pipeline Details", false)

	var sb strings.Builder

	status := pipeline.State.Result.Name
	icon := util.GetIconForStatus(status)
	statusColor := util.HexColor(util.GetColorForStatus(status))

	sb.WriteString(fmt.Sprintf("[::b]%s Pipeline #[-]%d [::b](Run #%d)[-]\n", constants.ICON_BUILD, pipeline.BuildNumber, pipeline.RunNumber))
	sb.WriteString(fmt.Sprintf("[::b]Status      :[-] [%s]%s[-] %s\n", statusColor, icon, status))
	sb.WriteString(fmt.Sprintf("[::b]Started     :[-] %s\n", util.FormatTime(pipeline.CreatedOn)))
	sb.WriteString(fmt.Sprintf("[::b]Completed   :[-] %s\n", util.FormatTime(pipeline.CompletedOn)))
	sb.WriteString(fmt.Sprintf("[::b]Duration    :[-] %d seconds\n", pipeline.Duration))

	sb.WriteString(fmt.Sprintf("\n[::b]%s Triggered by:[-] ", constants.ICON_COMMIT))
	if pipeline.Trigger.Type != "" {
		sb.WriteString(fmt.Sprintf("%s (%s)\n", pipeline.Trigger.Name, pipeline.Trigger.Type))
	} else {
		sb.WriteString("Unknown\n")
	}
	sb.WriteString(fmt.Sprintf("[::b]Created by :[-] %s\n", pipeline.Creator.DisplayName))

	sb.WriteString(fmt.Sprintf("\n[::b]%s Target[-]\n", constants.ICON_BRANCH))
	sb.WriteString(fmt.Sprintf("[::b]Type       :[-] %s\n", pipeline.Target.RefType))
	sb.WriteString(fmt.Sprintf("[::b]Name       :[-] %s\n", pipeline.Target.RefName))
	sb.WriteString(fmt.Sprintf("[::b]Commit     :[-] %s\n", pipeline.Target.Commit.Hash))

	textView.
		SetText(sb.String()).SetTextColor(tcell.ColorDarkGray)

	return textView
}
