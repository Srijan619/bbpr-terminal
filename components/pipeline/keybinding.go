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
		state.PipelineUIState.PipelineListFlex, state.PipelineUIState.PipelineSteps, state.PipelineUIState.PipelineStepCommandsView,
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
