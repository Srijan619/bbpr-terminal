package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

// GenerateStepDetailTable returns a tview.Table showing pipeline step details.

func GenerateStepCards(steps []types.StepDetail) *tview.Flex {
	mainFlexWrapper := tview.NewFlex()

	for _, step := range steps {
		var statusColor string
		switch step.State.Name {
		case "COMPLETED":
			statusColor = "[green]"
		case "FAILED":
			statusColor = "[red]"
		case "IN_PROGRESS":
			statusColor = "[yellow]"
		default:
			statusColor = "[gray]"
		}

		stepContent := fmt.Sprintf(
			"[::b]Step:[-:-] %s\n[::b]Status:[-:-] %s%s[-]\n[::b]Started:[-:-] %s\n[::b]Duration:[-:-] %ds\n",
			step.Name,
			statusColor, step.State.Name,
			step.StartedOn,
			step.DurationInSeconds,
		)

		stepView := tview.NewTextView().
			SetText(stepContent).
			SetDynamicColors(true).
			SetWrap(true).
			SetTextAlign(tview.AlignLeft).
			SetBorder(true).
			SetBackgroundColor(tcell.ColorDefault).
			SetTitle(fmt.Sprintf(" UUID: %s ", step.UUID))

		mainFlexWrapper.AddItem(stepView, 0, 1, false)
	}

	mainFlexWrapper.SetDirection(tview.FlexRow).
		SetBackgroundColor(tcell.ColorDefault)
	return mainFlexWrapper
}
