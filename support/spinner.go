package support

import (
	"simple-git-terminal/state"
	"time"

	"github.com/rivo/tview"
)

// ShowLoadingSpinner displays a loading spinner in the provided view while performing an async operation.
func ShowLoadingSpinner(view interface{}, fetch func() (interface{}, error), onComplete func(result interface{}, err error)) {
	ShowLoadingSpinnerWithApp(state.GlobalState.App, view, fetch, onComplete)
}

func ShowPipelineLoadingSpinner(view interface{}, fetch func() (interface{}, error), onComplete func(result interface{}, err error)) {
	ShowLoadingSpinnerWithApp(state.PipelineUIState.App, view, fetch, onComplete)
}

func ShowLoadingSpinnerWithApp(app *tview.Application, view interface{}, fetch func() (interface{}, error), onComplete func(result interface{}, err error)) {
	// Initial loading message
	UpdateView(view, "⠋ Loading...")

	go func() {
		spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		ticker := time.NewTicker(100 * time.Millisecond)
		done := make(chan bool)

		go func() {
			i := 0
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					app.QueueUpdateDraw(func() {
						UpdateView(view, spinChars[i]+" Loading...")
					})
					i = (i + 1) % len(spinChars)
				}
			}
		}()

		result, err := fetch()

		ticker.Stop()
		done <- true

		app.QueueUpdateDraw(func() {
			UpdateView(view, "")
			onComplete(result, err)
		})
	}()
}
