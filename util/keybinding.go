package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"simple-git-terminal/apis/bitbucket"

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
				log.Printf("[SAUSAGE ESC]...%d", currentFocusIndex)
				state.SetIsSearchMode(false)
				state.GlobalState.App.SetFocus(state.GlobalState.PrList) // Focus back to PrList or another view
			case tcell.KeyEnter:
				currentFocusIndex = 0
				state.SetSearchTerm(state.GlobalState.PrListSearchBar.GetText())
				state.SetIsSearchMode(false)
				state.GlobalState.PrList.Clear()
				go func() {
					// Perform search or exit search mode
					queryFetchedPRs := bitbucket.FetchPRsByQuery(bitbucket.BuildQuery(state.SearchTerm))
					state.SetFilteredPRs(&queryFetchedPRs)
					UpdatePRListView()
					state.GlobalState.App.SetFocus(state.GlobalState.PrList) // Focus back to PrList or another view
				}()
			default:
				return event // Ignore other keys in search mode
			}
			UpdateFocusBorders(focusOrder, currentFocusIndex, VIEW_ACTIVE_BORDER_COLOR)

		} else {
			// Handle keybindings when not in search mode
			switch event.Key() {
			case tcell.KeyTAB:
				// Cycle focus between views
				currentFocusIndex = (currentFocusIndex + 1) % len(focusOrder)
				if currentFocusIndex >= len(focusOrder) {
					currentFocusIndex = 0 // If we go out of bounds, set to the first element
				}
				log.Printf("[SAUSAGE TAB]...%d", currentFocusIndex)
				state.GlobalState.App.SetFocus(focusOrder[currentFocusIndex])

			case tcell.KeyCtrlC:
				state.GlobalState.App.Stop()

			case tcell.KeyRune:
				switch event.Rune() {
				case 's':
					// Search mode
					state.SetIsSearchMode(true)
					currentFocusIndex = len(focusOrder) - 1
					log.Printf("[SAUSAGE S]...%d", currentFocusIndex)

					state.GlobalState.App.SetFocus(state.GlobalState.PrListSearchBar)
					go func() {
						state.GlobalState.PrListSearchBar.SetText("", true)
					}()
				case 't', 'T':
					// Focus on DiffStatView (T or t)
					currentFocusIndex = len(focusOrder) - 3
					state.GlobalState.App.SetFocus(state.GlobalState.DiffStatView)

				case 'c', 'C':
					// Focus on DiffDetails (C or c)
					currentFocusIndex = len(focusOrder) - 2
					state.GlobalState.App.SetFocus(state.GlobalState.DiffDetails)

				case 'a', 'A':
					// Focus on ActivityView (A or a)
					currentFocusIndex = len(focusOrder) - 4
					state.GlobalState.App.SetFocus(state.GlobalState.ActivityView)

				case 'p':
					// Focus on PR List
					currentFocusIndex = 0
					state.GlobalState.App.SetFocus(state.GlobalState.PrList)

				case 'd', 'D':
					// Focus on PR Details (D or d)
					currentFocusIndex = len(focusOrder) - 5
					state.GlobalState.App.SetFocus(state.GlobalState.PrDetails)

				case 'q':
					// Quit application
					state.GlobalState.App.SetRoot(state.GlobalState.MainFlexWrapper, true)

				case 'm', 'o', 'r':
					// Toggle PR filters
					switch event.Rune() {
					case 'm':
						state.SetPRStatusFilter("merged", !state.PRStatusFilter.Merged)
					case 'o':
						state.SetPRStatusFilter("open", !state.PRStatusFilter.Open)
					case 'r':
						state.SetPRStatusFilter("declined", !state.PRStatusFilter.Declined)
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
