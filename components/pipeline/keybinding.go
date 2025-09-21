package pipeline

import (
	"simple-git-terminal/state"
	"simple-git-terminal/support"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	VIEW_ACTIVE_BORDER_COLOR = tcell.ColorOrange
)

func SetupKeyBindings() {
	focusOrder := []tview.Primitive{
		state.PipelineUIState.PipelineList, state.PipelineUIState.PipelineSteps, state.PipelineUIState.PipelineStepCommandsView,
	}
	// Define focus order
	currentFocusIndex := 0
	support.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

	state.PipelineUIState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle keybindings when not in search mode
		switch event.Key() {
		case tcell.KeyTAB:
			// Cycle focus between views
			currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
			if currentFocusIndex >= len(focusOrder) {
				currentFocusIndex = 0 // If we go out of bounds, set to the first element
			}
			state.PipelineUIState.App.SetFocus(focusOrder[currentFocusIndex])

			support.SetTableSelectability(focusOrder, currentFocusIndex, map[tview.Primitive]tview.Primitive{
				state.PipelineUIState.PipelineSteps:            state.PipelineUIState.PipelineStepsTable,
				state.PipelineUIState.PipelineStepCommandsView: state.PipelineUIState.PipelineScriptCommandsTable,
			})

		case tcell.KeyRune:
			switch event.Rune() {
			case 'r':
				PopulatePipelineList()
			}
		}
		// Update focus borders after focus change
		support.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

		return event
	})
}
