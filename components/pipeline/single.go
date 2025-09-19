package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"strings"
)

func GeneratePPDetail(steps []types.StepDetail) string {
	var sb strings.Builder
	for _, step := range steps {
		sb.WriteString(fmt.Sprintf(
			"Step UUID: %s\nName: %s\nState: %s\nStarted On: %s\nDuration: %d seconds\n\n",
			step.UUID,
			step.Name,
			step.State.Name,
			step.StartedOn,
			step.DurationInSeconds,
		))
	}
	return sb.String()
}
