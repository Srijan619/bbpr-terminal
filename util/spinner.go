package util

import (
	"time"

	"github.com/rivo/tview"

	"simple-git-terminal/state"
)

// ShowLoadingSpinner displays a loading spinner in the provided view while performing an async operation.
func ShowLoadingSpinner(view *tview.Flex, fetch func() (string, error), onComplete func(result string, err error)) {
	// Initial loading message
	UpdateView(state.GlobalState.App, view, "⠋ Loading...")

	// Run the fetch operation in a goroutine
	go func() {
		// Simulate spinner animation
		spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		ticker := time.NewTicker(100 * time.Millisecond)
		done := make(chan bool)

		// Spinner animation loop
		go func() {
			i := 0
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					state.GlobalState.App.QueueUpdateDraw(func() {
						UpdateView(state.GlobalState.App, view, spinChars[i]+" Loading...")
					})
					i = (i + 1) % len(spinChars)
				}
			}
		}()

		// Perform the async operation
		result, err := fetch()

		// Stop spinner animation
		ticker.Stop()
		done <- true

		// Update the view with the result (back on the main thread)
		state.GlobalState.App.QueueUpdateDraw(func() {
			onComplete(result, err)
		})
	}()
}

