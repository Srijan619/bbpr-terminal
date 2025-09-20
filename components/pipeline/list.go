package pipeline

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/ui"
	"simple-git-terminal/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func PopulatePipelineList(ppTable *tview.Table) *tview.Table {
	ppTable.SetBackgroundColor(tcell.ColorDefault)

	var (
		pipelineList  []types.PipelineResponse
		nextPageURL   string
		lastFetchDone bool
		isLoading     bool
	)

	loadPipelines := func(query string, appendData bool) {
		if lastFetchDone {
			log.Println("[INFO] Already fetched this page, skipping...")
			return
		}

		if isLoading {
			return
		}
		isLoading = true

		util.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineList, func() (interface{}, error) {
			pps, pagination := bitbucket.FetchPipelinesByQuery(query)
			if pps == nil {
				log.Println("Failed to fetch pipelines, nil returned")
				return nil, fmt.Errorf("failed to fetch pipelines")
			}

			// If the number of pipelines fetched is less than the batch (max of 10),
			// this means we have reached the last page and there are no more pipelines to fetch.
			// Therefore, clear nextPageURL to prevent further pagination.
			//
			// This is based on the assumption that the query contains a 'pagelen' parameter indicating the page size.
			if pagination.PageLen < 10 {
				log.Printf("[INFO] No more pipelines to fetch, reached less than 10 page leng. %d", len(pps))
				lastFetchDone = true
			} else {
				nextPageURL = pagination.Next
			}

			nextPageURL = pagination.Next
			return pps, nil
		}, func(result interface{}, err error) {
			defer func() { isLoading = false }()

			pps, ok := result.([]types.PipelineResponse)
			if !ok {
				util.UpdateView(state.PipelineUIState.PipelineList, fmt.Sprintf("[red]Error: %v[-]", err))
				return
			}

			if appendData {
				pipelineList = append(pipelineList, pps...)
			} else {
				pipelineList = pps
			}

			ui.PopulatePPList(ppTable, pipelineList)

			ppTable.SetSelectedFunc(func(row, column int) {
				go func() {
					HandleOnPipelineSelect(pipelineList, row)
				}()
			})
		})
	}

	// Initial fetch
	go loadPipelines(bitbucket.BuildQuery(""), false)

	// Handle scroll near bottom to fetch more
	ppTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		row, _ := ppTable.GetSelection()
		totalRows := ppTable.GetRowCount()

		if key == tcell.KeyDown && row >= totalRows-2 && nextPageURL != "" {
			log.Printf("[INFO] Scrolling near bottom, loading next page: %s", nextPageURL)

			// Extract query string only
			query := util.ExtractQueryFromNextURL(nextPageURL)

			// Avoid firing multiple times
			nextPageURL = ""

			// Load next page in background
			go loadPipelines(query, true)
		}

		return event
	})
	return ppTable
}

func HandleOnPipelineSelect(pipelines []types.PipelineResponse, row int) {
	// Validate row index
	if row < 0 || row >= len(pipelines) {
		log.Printf("Invalid row index: %d, pipelines count: %d", row, len(pipelines))
		return
	}

	if state.PipelineUIState == nil {
		log.Println("PipelineUIState is nil, cannot update UI")
		return
	}
	if state.PipelineUIState.PipelineSteps == nil {
		log.Println("PipelineSteps view is nil, cannot update UI")
		return
	}

	selectedPipeline := pipelines[row]
	log.Printf("Selected pipeline UUID: %s, Name: %s", selectedPipeline.UUID, selectedPipeline.Creator.DisplayName)

	util.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineSteps, func() (interface{}, error) {
		EmptyAllPipelineListDependentViews()

		steps := bitbucket.FetchPipelineSteps(selectedPipeline.UUID)
		if steps == nil {
			log.Println("Failed to fetch pipeline steps, nil returned")
			return nil, fmt.Errorf("failed to fetch pipeline steps")
		}
		return steps, nil
	}, func(result interface{}, err error) {
		steps, ok := result.([]types.StepDetail)
		if !ok {
			util.UpdateView(state.PipelineUIState.PipelineSteps, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		view := GenerateStepsView(steps, selectedPipeline)

		util.UpdateView(state.PipelineUIState.PipelineStepsDebugView, GeneratePPDebugInfo(selectedPipeline))
		util.UpdateView(state.PipelineUIState.PipelineSteps, view)

		HandleOnStepSelect(steps, selectedPipeline, 0) // Auto select first step and fetch the info
	})
}

func EmptyAllPipelineListDependentViews() {
	// Always start a fresh slate for rest of dependent views
	emptyView := tview.NewFlex()
	util.UpdateView(state.PipelineUIState.PipelineStepsDebugView, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineSteps, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineStep, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineStepCommandLogView, emptyView)
}
