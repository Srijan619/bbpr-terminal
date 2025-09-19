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

	// Group steps by parallel group name
	grouped := make(map[string][]types.StepDetail)
	ungrouped := []types.StepDetail{}

	for _, step := range steps {
		// This is not yet supported by open API but internal API has this
		if step.ParallelGroup.GroupName != "" {
			grouped[step.ParallelGroup.GroupName] = append(grouped[step.ParallelGroup.GroupName], step)
		} else {
			ungrouped = append(ungrouped, step)
		}
	}

	// Add grouped steps
	for groupName, groupSteps := range grouped {
		groupNode := createNode(fmt.Sprintf(" Parallel Group: %s", groupName), tcell.ColorLightBlue).SetExpanded(true)

		for _, step := range groupSteps {
			stepNode := buildStepNode(step, createNode)
			groupNode.AddChild(stepNode)
		}

		rootNode.AddChild(groupNode)
	}

	// Add ungrouped steps
	for _, step := range ungrouped {
		stepNode := buildStepNode(step, createNode)
		rootNode.AddChild(stepNode)
	}

	tree := tview.NewTreeView().
		SetRoot(rootNode).
		SetCurrentNode(rootNode).
		SetGraphics(true)

	tree.SetBackgroundColor(tcell.ColorDefault)

	return tree
}

func buildStepNode(step types.StepDetail, createNode func(string, tcell.Color) *tview.TreeNode) *tview.TreeNode {
	icon := util.GetIconForStatus(step.State.Result.Name)
	color := util.GetColorForStatus(step.State.Result.Name)

	status := fmt.Sprintf("[#%s]%s[-]", util.HexColor(color), icon)

	title := fmt.Sprintf(" %s [::b]%s[-]", status, step.Name)
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

	return stepNode
}
