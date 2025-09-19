package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"strings"
	"time"
)

func GeneratePPDebugInfo(pipeline *types.PipelineResponse) string {
	if pipeline == nil {
		return "No pipeline data available."
	}

	formatTime := func(t string) string {
		parsed, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return t
		}
		return parsed.Format("2006-01-02 15:04:05")
	}

	getStatusEmoji := func(status types.PipelineStatus) string {
		switch {
		case status.Successful():
			return "✅"
		case status.Failed():
			return "❌"
		case status.Running():
			return "🏃"
		case status.Pending():
			return "⏳"
		case status.Stopped():
			return "⛔"
		default:
			return "❔"
		}
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🛠️  Pipeline #%d (Run #%d)\n", pipeline.BuildNumber, pipeline.RunNumber))
	sb.WriteString(fmt.Sprintf("Status      : %s %s\n", pipeline.State.Result.Name, getStatusEmoji(pipeline.State.Result.Name)))
	sb.WriteString(fmt.Sprintf("Started     : %s\n", formatTime(pipeline.CreatedOn)))
	sb.WriteString(fmt.Sprintf("Completed   : %s\n", formatTime(pipeline.CompletedOn)))
	sb.WriteString(fmt.Sprintf("Duration    : %d seconds\n", pipeline.Duration))

	sb.WriteString("\n👤 Triggered by: ")
	if pipeline.Trigger.Type != "" {
		sb.WriteString(fmt.Sprintf("%s (%s)\n", pipeline.Trigger.Name, pipeline.Trigger.Type))
	} else {
		sb.WriteString("Unknown\n")
	}
	sb.WriteString(fmt.Sprintf("Created by : %s\n", pipeline.Creator.DisplayName))

	sb.WriteString("\n🔀 Target\n")
	sb.WriteString(fmt.Sprintf("Type       : %s\n", pipeline.Target.RefType))
	sb.WriteString(fmt.Sprintf("Name       : %s\n", pipeline.Target.RefName))
	sb.WriteString(fmt.Sprintf("Commit     : %s\n", pipeline.Target.Commit.Hash))

	return sb.String()
}
