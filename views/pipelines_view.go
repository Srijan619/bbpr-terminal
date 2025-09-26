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

type PipelinesView struct {
	BaseView
	table *widgets.PipelineTable
	bus   *events.Bus
}

func NewPipelineView(bus *events.Bus) *PipelinesView {
	table := widgets.NewPipelineTable()
	sv := &PipelinesView{
		BaseView: NewBaseView(table), // gives Subscribe, Refresh, etc.
		table:    table,
		bus:      bus,
	}

	sv.Subscribe(bus) // attach reactive event handling

	sv.table.SetSelectedFunc(func(row, column int) {
		sv.table.SetSelectedRow(row) // <- this updates SelectedRow and returns pipeline

		selected := sv.table.GetSelectedPipeline()
		if selected == nil {
			log.Printf("[ERROR] No pipeline at row %d", row)
			return
		}

		// Now publish the selection event
		sv.bus.Publish(events.PipelineSelectedEvent{
			Pipeline: *selected,
			Row:      row,
		})
	})

	sv.Render()

	return sv
}

// Subscribe handles events for this view
func (sv *PipelinesView) Subscribe(bus *events.Bus) {
	bus.Subscribe(func(e events.Event) {
	})
}

func (sv *PipelinesView) Render() {
	support.ShowPipelineLoadingSpinner(sv.table, func() (interface{}, error) {
		pps, _ := bitbucket.FetchPipelinesByQuery("")
		if pps == nil {
			log.Println("Failed to fetch pipelines, nil returned")
			return nil, fmt.Errorf("failed to fetch pipelines")
		}

		return pps, nil
	}, func(result interface{}, err error) {
		pps, ok := result.([]types.PipelineResponse)
		if !ok {
			support.UpdateView(sv.table, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}
		sv.table.SetPipelines(pps, 0)
	})
}

func (sv *PipelinesView) GetView() *widgets.PipelineTable {
	return sv.table
}

// Refresh triggers a full refresh (from BaseView)
func (sv *PipelinesView) Refresh() {
	sv.table.Clear()
	sv.Render()
}
