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

func GenerateStepTreeView(steps []types.StepDetail, selectedPipeline types.PipelineResponse) *tview.TreeView {
	createNode := func(text string, color tcell.Color) *tview.TreeNode {
		return tview.NewTreeNode(text).
			SetColor(color)
	}

	rootNode := createNode(fmt.Sprintf("Build %d", selectedPipeline.BuildNumber), tcell.ColorDarkGray)

	for _, step := range steps {
		var statusIcon string
		switch step.State.Name {
		case "COMPLETED":
			statusIcon = "[green]✔[-]"
		case "FAILED":
			statusIcon = "[red]✖[-]"
		case "IN_PROGRESS":
			statusIcon = "[yellow]…[-]"
		default:
			statusIcon = "[gray]?[-]"
		}

		title := fmt.Sprintf("%s [::b]%s[-] %s", statusIcon, step.Name, step.State.Name)
		stepNode := createNode(title, tcell.ColorWhite).
			SetReference(step).
			SetExpanded(false)

		// Add step info as child nodes
		stepNode.AddChild(createNode(fmt.Sprintf("UUID: %s", step.UUID), tcell.ColorLightGrey))
		stepNode.AddChild(createNode(fmt.Sprintf("Duration: %ds", step.DurationInSeconds), tcell.ColorLightGrey))
		stepNode.AddChild(createNode(fmt.Sprintf("Started: %s", step.StartedOn), tcell.ColorLightGrey))
		stepNode.AddChild(createNode(fmt.Sprintf("Completed: %s", step.CompletedOn), tcell.ColorLightGrey))

		// Commands as child group
		if len(step.ScriptCommands) > 0 {
			commandsNode := createNode("Commands:", tcell.ColorLightBlue)
			for _, cmd := range step.ScriptCommands {
				commandsNode.AddChild(createNode(fmt.Sprintf("- %s", cmd.Command), tcell.ColorGrey))
			}
			stepNode.AddChild(commandsNode)
		}

		rootNode.AddChild(stepNode)
	}

	tree := tview.NewTreeView().
		SetRoot(rootNode).
		SetCurrentNode(rootNode).
		SetGraphics(true)

	tree.
		SetBackgroundColor(tcell.ColorDefault)

	return tree
}
