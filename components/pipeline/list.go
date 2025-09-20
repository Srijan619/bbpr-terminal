package pipeline

import (
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/ui"
	"simple-git-terminal/util"
	"time"

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
		frame         int
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

			ui.PopulatePPList(ppTable, pipelineList, frame)

			ppTable.SetSelectedFunc(func(row, column int) {
				go func() {
					HandleOnPipelineSelect(pipelineList, row)
				}()
			})
		})
	}

	// Initial fetch
	go loadPipelines(bitbucket.BuildQuery(""), false)

	// Animate pipeline list with frame counter
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			if len(pipelineList) == 0 {
				continue
			}
			frame++

			state.PipelineUIState.App.QueueUpdateDraw(func() {
				for i, pp := range pipelineList {
					status := pp.State.Result.Name
					if status == "" {
						status = pp.State.Name
					}

					if !(status.InProgress() || status.Running() || status.Pending()) {
						continue
					}

					// Fetch new state
					updatedPipeline := bitbucket.FetchPipeline(pp.UUID)

					icon := util.GetIconForStatusWithColorAnimated(updatedPipeline.State.Name, frame)
					color := util.GetColorForStatus(updatedPipeline.State.Name)
					text := fmt.Sprintf("%s %s", icon, updatedPipeline.State.Name)

					ppTable.SetCell(i, 6, util.CellFormat(text, color)) // Column 6 is status
				}
			})
		}
	}()

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

		if selectedPipeline.State.Name.InProgress() {
			// if selected pipeline is in progress, track it with constant polling in backgroun until it is completed
			go TrackPipelineLive(selectedPipeline)
		}
		view := GenerateStepsView(steps, selectedPipeline)

		util.UpdateView(state.PipelineUIState.PipelineStepsDebugView, GeneratePPDebugInfo(selectedPipeline))
		util.UpdateView(state.PipelineUIState.PipelineSteps, view)

		HandleOnStepSelect(steps, selectedPipeline, 0) // Auto select first step and fetch the info
	})
}

func TrackPipelineLive(pipeline types.PipelineResponse) {
	log.Printf("[INFO] Starting live tracking for pipeline: %s", pipeline.UUID)

	ticker := time.NewTicker(3 * time.Second) // poll every 3 seconds
	defer ticker.Stop()

	for range ticker.C {
		updatedPipeline := bitbucket.FetchPipeline(pipeline.UUID)
		if updatedPipeline.UUID == "" {
			log.Printf("[WARN] Could not fetch pipeline update for %s", pipeline.UUID)
			continue
		}

		if !updatedPipeline.State.Name.InProgress() {
			log.Printf("[INFO] Pipeline %s is no longer in progress. Stopping tracker.", pipeline.UUID)
			break
		}

		steps := bitbucket.FetchPipelineSteps(updatedPipeline.UUID)
		if steps == nil {
			log.Printf("[WARN] Could not fetch updated steps for pipeline %s", updatedPipeline.UUID)
			continue
		}

		// Re-render step view
		view := GenerateStepsView(steps, *updatedPipeline)
		util.UpdateView(state.PipelineUIState.PipelineSteps, view)
	}
}

func EmptyAllPipelineListDependentViews() {
	// Always start a fresh slate for rest of dependent views
	emptyView := tview.NewFlex()
	util.UpdateView(state.PipelineUIState.PipelineStepsDebugView, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineSteps, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineStep, emptyView)
	util.UpdateView(state.PipelineUIState.PipelineStepCommandLogView, emptyView)
}
