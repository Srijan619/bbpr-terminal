package pipeline

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GenerateStepView(step types.StepDetail, selectedPipeline types.PipelineResponse) tview.Primitive {
	// ─── TEXT VIEW: Metadata ─────────────────────────────────────────────
	textView := util.CreateTextviewComponent("Step Details", false)

	var sb strings.Builder

	// Status + color/icon
	status := step.State.Name
	icon := util.GetIconForStatus(status)
	color := util.GetColorForStatus(status)
	colorHex := util.HexColor(color)

	sb.WriteString(fmt.Sprintf("[::b]UUID       :[-] %s\n", step.UUID))
	sb.WriteString(fmt.Sprintf("[::b]Status     :[-] [#%s]%s[-] %s\n", colorHex, icon, status))
	sb.WriteString(fmt.Sprintf("[::b]Started    :[-] %s\n", util.FormatTime(step.StartedOn)))
	sb.WriteString(fmt.Sprintf("[::b]Completed  :[-] %s\n", util.FormatTime(step.CompletedOn)))

	// Setup commands (optional)
	if len(step.SetupCommands) > 0 {
		sb.WriteString("\n[::b]🔧 Setup Commands:[-]\n")
		for _, cmd := range step.SetupCommands {
			sb.WriteString(fmt.Sprintf("  %s [::b]%s\n", util.ICON_SIDE_ARROW, cmd.Name))
		}
	}

	textView.SetText(sb.String()).SetTextColor(tcell.ColorDarkGray)

	// ─── TABLE VIEW: Script Commands ─────────────────────────────────────
	scriptTable := tview.NewTable()

	scriptTable.
		SetBorders(false).
		SetSelectable(true, false).
		SetBackgroundColor(tcell.ColorDefault)

	scriptTable.SetTitle("Script Commands").SetTitleAlign(tview.AlignLeft).SetBorder(true)

	for i, cmd := range step.ScriptCommands {
		cmdText := fmt.Sprintf("%s [::b]%s", util.ICON_SIDE_ARROW, cmd.Name)
		scriptTable.SetCell(i, 0, util.CellFormat(cmdText, tcell.ColorWhite))
	}

	if len(step.ScriptCommands) == 0 {
		scriptTable.SetCell(0, 0, util.CellFormat(" No script commands available", tcell.ColorGray))
	}

	// ─── MAIN FLEX ───────────────────────────────────────────────────────
	layout := tview.NewFlex()

	layout.
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 3, false).
		AddItem(scriptTable, 0, 1, true)

	scriptTable.SetSelectedFunc(func(row, column int) {
		go func() {
			HandleOnScriptCommandSelected(step.SetupCommands, step, selectedPipeline, row)
		}()
	})
	scriptTable.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkOrange))
	return layout
}

func HandleOnScriptCommandSelected(commands []types.CommandDetail, selectedStep types.StepDetail, selectedPipeline types.PipelineResponse, row int) {
	// Validate row index
	if row < 0 || row >= len(commands) {
		log.Printf("Invalid row index: %d, commands count: %d", row, len(commands))
		return
	}

	if state.PipelineUIState == nil {
		log.Println("PipelineUIState is nil, cannot update UI")
		return
	}

	selectedCommand := commands[row]

	util.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineStepCommandLogView, func() (interface{}, error) {
		commandLog, error := bitbucket.FetchPipelineStepCommandLogs(selectedPipeline.UUID, selectedStep.UUID, selectedCommand.Name)

		if error != nil {
			return nil, fmt.Errorf("failed to fetch single step command %s", selectedCommand.Name)
		}

		return commandLog, nil
	}, func(result interface{}, err error) {
		commandLog, ok := result.(string)
		if !ok {
			util.UpdateView(state.PipelineUIState.PipelineStepCommandLogView, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		util.UpdateView(state.PipelineUIState.PipelineStepCommandLogView, GenerateStepCommandLogView(commandLog, selectedCommand.Name))
	})
}
