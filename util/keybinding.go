package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/state"
)

const (
	VIEW_ACTIVE_BORDER_COLOR = tcell.ColorOrange
)

func SetupKeyBindings(callback func()) {
	focusOrder := []tview.Primitive{
		state.GlobalState.PrListFlex, state.GlobalState.PrDetails, state.GlobalState.ActivityView,
		state.GlobalState.DiffStatView, state.GlobalState.DiffDetails, state.GlobalState.PrListSearchBar,
	}
	// Define focus order
	currentFocusIndex := 0
	UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

	state.GlobalState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If in search mode, only allow Esc or Enter keys
		if state.IsSearchMode {
			switch event.Key() {
			case tcell.KeyEsc:
				currentFocusIndex = 0
				state.SetIsSearchMode(false)
				state.GlobalState.App.SetFocus(state.GlobalState.PrList) // Focus back to PrList or another view
				log.Printf("Esc pressed escaping now......")
				UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
			case tcell.KeyEnter:
				currentFocusIndex = 0
				state.SetSearchTerm(state.GlobalState.PrListSearchBar.GetText())
				ShowSpinnerFetchPRsByQueryAndUpdatePrList()
			default:
				return event // Ignore other keys in search mode
			}
		} else {
			// Handle keybindings when not in search mode
			switch event.Key() {
			case tcell.KeyTAB:
				// Cycle focus between views
				currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
				if currentFocusIndex >= len(focusOrder) {
					currentFocusIndex = 0 // If we go out of bounds, set to the first element
				}
				state.GlobalState.App.SetFocus(focusOrder[currentFocusIndex])

			case tcell.KeyCtrlC:
				state.GlobalState.App.Stop()

			case tcell.KeyRune:
				switch event.Rune() {
				case 's':
					// Search mode
					state.SetIsSearchMode(true)
					currentFocusIndex = len(focusOrder) - 1

					state.GlobalState.App.SetFocus(state.GlobalState.PrListSearchBar)
					//state.GlobalState.PrListSearchBar.SetText("")
					UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR) // TODO: This is repeated here as we need to return nil from event rune otherwise it adds pressed key rune to textarea
					return nil

				case 't', 'T':
					currentFocusIndex = len(focusOrder) - 3
					switch event.Rune() {
					case 't':
						state.GlobalState.App.SetFocus(state.GlobalState.DiffStatView)
					case 'T':
						state.GlobalState.App.SetRoot(state.GlobalState.DiffStatView, true)
					}

				case 'c', 'C':
					currentFocusIndex = len(focusOrder) - 2
					switch event.Rune() {
					case 'c':
						state.GlobalState.App.SetFocus(state.GlobalState.DiffDetails)
					case 'C':
						state.GlobalState.App.SetRoot(state.GlobalState.DiffDetails, true)
					}

				case 'a', 'A':
					currentFocusIndex = len(focusOrder) - 4
					switch event.Rune() {
					case 'a':
						state.GlobalState.App.SetFocus(state.GlobalState.ActivityView)
					case 'A':
						state.GlobalState.App.SetRoot(state.GlobalState.ActivityView, true)
					}

				case 'p', 'P':
					currentFocusIndex = 0
					switch event.Rune() {
					case 'p':
						state.GlobalState.App.SetFocus(state.GlobalState.PrList)
					case 'P':
						state.GlobalState.App.SetRoot(state.GlobalState.PrList, true)
					}

				case 'd', 'D':
					currentFocusIndex = len(focusOrder) - 5
					switch event.Rune() {
					case 'd':
						state.GlobalState.App.SetFocus(state.GlobalState.PrDetails)
					case 'D':
						state.GlobalState.App.SetRoot(state.GlobalState.PrDetails, true)
					}

				case 'q':
					currentFocusIndex = 0
					state.GlobalState.App.SetRoot(state.GlobalState.MainFlexWrapper, true)

				case 'm', 'o', 'r', 'i':
					// Toggle PR filters
					switch event.Rune() {
					case 'm':
						UpdatePRListWithFilter("merged", !state.PRStatusFilter.Merged)
					case 'o':
						UpdatePRListWithFilter("open", !state.PRStatusFilter.Open)
					case 'r':
						UpdatePRListWithFilter("declined", !state.PRStatusFilter.Declined)
					case 'i':
						UpdatePRListWithFilter("iamreviewing", !state.PRStatusFilter.IAmReviewing)
					}
					callback()
				}
			}
			// Update focus borders after focus change
			UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)
		}

		return event
	})
}
