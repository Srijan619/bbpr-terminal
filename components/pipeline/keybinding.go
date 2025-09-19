package pipeline

import (
	"simple-git-terminal/state"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	VIEW_ACTIVE_BORDER_COLOR = tcell.ColorOrange
)

func SetupKeyBindings() {
	focusOrder := []tview.Primitive{
		state.PipelineUIState.PipelineListFlex, state.PipelineUIState.PipelineSteps, state.PipelineUIState.PipelineStep,
	}
	// Define focus order
	currentFocusIndex := 0
	util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

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
		}
		// Update focus borders after focus change
		util.UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

		return event
	})
}
