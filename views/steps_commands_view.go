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

// StepsCommandsView is a reactive view that renders a StepsTable and manages events
type StepsCommandsView struct {
	BaseView
	table *widgets.StepsCommandsTable
	bus   *events.Bus
}

func NewStepsCommandsView(bus *events.Bus) *StepsCommandsView {
	table := widgets.NewStepsCommandsTable()
	sv := &StepsCommandsView{
		BaseView: NewBaseView(table), // gives Subscribe, Refresh, etc.
		table:    table,
		bus:      bus,
	}

	sv.Subscribe(bus) // attach reactive event handling

	return sv
}

// Subscribe handles events for this view
func (sv *StepsCommandsView) Subscribe(bus *events.Bus) {
	bus.Subscribe(func(e events.Event) {
		switch ev := e.(type) {
		case events.StepSelectedEvent:
			support.ShowPipelineLoadingSpinner(sv.table, func() (interface{}, error) {
				step := bitbucket.FetchPipelineStep(ev.Step.Pipeline.UUID, ev.Step.UUID)
				if step.UUID == "" {
					log.Println("Failed to fetch selected step's step, nil returned")
					return nil, fmt.Errorf("failed to fetch selected step's step")
				}
				return step, nil
			}, func(result interface{}, err error) {
				step, ok := result.(types.StepDetail)
				if !ok {
					support.UpdateView(sv.table, fmt.Sprintf("[red]Error: %v[-]", err))
					return
				}

				sv.table.SetStepDetail(step, 0)
			})
		case events.StepsUpdatedEvent:
			// if ev.PipelineUUID == sv.table.GetSelectedStep {
			// 	sv.table.RenderSteps(ev.Steps) // let table widget handle actual cell updates
			// }
		}
	})
}

// Render returns the table widget as tview.Primitive
func (sv *StepsCommandsView) Render() *widgets.StepsCommandsTable {
	return sv.table
}

// Refresh triggers a full refresh (from BaseView)
func (sv *StepsCommandsView) Refresh() {
	sv.table.Clear()
	// if sv.OnRefresh != nil {
	// 	sv.OnRefresh()
	// }
}
