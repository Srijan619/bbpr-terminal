package pipeline

import (
	"fmt"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
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

func GenerateStepTreeView(steps []types.StepDetail, selectedPipeline types.PipelineResponse) *tview.TreeView {
	createNode := func(text string, fg tcell.Color) *tview.TreeNode {
		style := tcell.StyleDefault.Foreground(fg).Background(tcell.ColorDefault)
		selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkCyan).Background(tcell.ColorDefault)

		return tview.NewTreeNode(text).
			SetTextStyle(style).
			SetSelectedTextStyle(selectedStyle)
	}

	rootNode := createNode(fmt.Sprintf("%s Build %d", util.ICON_BUILD, selectedPipeline.BuildNumber), tcell.ColorDefault)

	for _, step := range steps {
		iconWithColor := util.GetIconForStatusWithColor(step.State.Result.Name)

		title := fmt.Sprintf(" %s [::b]%s[-]", iconWithColor, step.Name)
		stepNode := createNode(title, tcell.ColorDefault).
			SetReference(step).
			SetExpanded(false)

		stepNode.AddChild(createNode(fmt.Sprintf("UUID: %s", step.UUID), tcell.ColorDefault))
		stepNode.AddChild(createNode(fmt.Sprintf("Duration: %ds", step.DurationInSeconds), tcell.ColorDefault))
		stepNode.AddChild(createNode(fmt.Sprintf("Started: %s", step.StartedOn), tcell.ColorDefault))
		stepNode.AddChild(createNode(fmt.Sprintf("Completed: %s", step.CompletedOn), tcell.ColorDefault))

		if len(step.ScriptCommands) > 0 {
			cmds := createNode("Commands:", tcell.ColorLightBlue)
			for _, cmd := range step.ScriptCommands {
				cmds.AddChild(createNode(fmt.Sprintf("- %s", cmd.Command), tcell.ColorDefault))
			}
			stepNode.AddChild(cmds)
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
