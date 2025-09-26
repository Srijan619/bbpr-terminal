package views

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/events"
	"simple-git-terminal/support"
	"simple-git-terminal/types"
	"simple-git-terminal/widgets"
)

// StepsView is a reactive view that renders a StepsTable and manages events
type StepsView struct {
	BaseView
	table *widgets.StepsTable
	bus   *events.Bus
}

func NewStepsView(bus *events.Bus) *StepsView {
	table := widgets.NewStepsTable()
	sv := &StepsView{
		BaseView: NewBaseView(table), // gives Subscribe, Refresh, etc.
		table:    table,
		bus:      bus,
	}

	sv.Subscribe(bus) // attach reactive event handling

	return sv
}

// Subscribe handles events for this view
func (sv *StepsView) Subscribe(bus *events.Bus) {
	bus.Subscribe(func(e events.Event) {
		switch ev := e.(type) {
		case events.PipelineSelectedEvent:
			// fetch steps async, then publish StepsUpdatedEvent
			support.ShowPipelineLoadingSpinner(sv.table, func() (interface{}, error) {
				steps := bitbucket.FetchPipelineSteps(ev.Pipeline.UUID)
				if steps == nil {
					log.Println("Failed to fetch pipeline steps, nil returned")
					return nil, fmt.Errorf("failed to fetch pipeline steps")
				}
				return steps, nil
			}, func(result interface{}, err error) {
				steps, ok := result.([]types.StepDetail)
				if !ok {
					support.UpdateView(sv.table, fmt.Sprintf("[red]Error: %v[-]", err))
					return
				}

				sv.table.SetSteps(steps, 0)
			})
		case events.StepsUpdatedEvent:
			// if ev.PipelineUUID == sv.table.GetSelectedStep {
			// 	sv.table.RenderSteps(ev.Steps) // let table widget handle actual cell updates
			// }
		}
	})
}

// Render returns the table widget as tview.Primitive
func (sv *StepsView) Render() *widgets.StepsTable {
	return sv.table
}

// Refresh triggers a full refresh (from BaseView)
func (sv *StepsView) Refresh() {
	sv.table.Clear()
	// if sv.OnRefresh != nil {
	// 	sv.OnRefresh()
	// }
}
