package pipeline

import (
	"context"
	"fmt"
	"log"
	"simple-git-terminal/apis/bitbucket"
	"simple-git-terminal/state"
	"simple-git-terminal/support"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func PopulatePipelineList() {
	var (
		pipelineList  []types.PipelineResponse
		nextPageURL   string
		lastFetchDone bool
		isLoading     bool
		frame         int
	)

	// Caching pipeline fetch results to reduce frequent API calls
	type pipelineCacheEntry struct {
		lastFetched time.Time
		data        types.PipelineResponse
	}

	pipelineCache := make(map[string]pipelineCacheEntry)

	loadPipelines := func(query string, appendData bool) {
		if lastFetchDone {
			log.Println("[INFO] Already fetched this page, skipping...")
			return
		}

		if isLoading {
			return
		}
		isLoading = true

		support.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineList, func() (interface{}, error) {
			pps, pagination := bitbucket.FetchPipelinesByQuery(query)
			if pps == nil {
				log.Println("Failed to fetch pipelines, nil returned")
				return nil, fmt.Errorf("failed to fetch pipelines")
			}

			// If page size is less than 10, we assume it's the last page
			if pagination.PageLen < 10 {
				log.Printf("[INFO] No more pipelines to fetch, less than page size. Fetched: %d", len(pps))
				lastFetchDone = true
			} else {
				nextPageURL = pagination.Next
			}

			return pps, nil
		}, func(result interface{}, err error) {
			defer func() { isLoading = false }()

			pps, ok := result.([]types.PipelineResponse)
			if !ok {
				support.UpdateView(state.PipelineUIState.PipelineList, fmt.Sprintf("[red]Error: %v[-]", err))
				return
			}

			if appendData {
				pipelineList = append(pipelineList, pps...)
			} else {
				pipelineList = pps
			}

			state.PipelineUIState.PipelineList.SetPipelines(pipelineList, frame)

			state.PipelineUIState.PipelineList.SetSelectedFunc(func(row, column int) {
				go func() {
					HandleOnPipelineSelect(pipelineList, row, frame)
				}()
			})

			// Select first pipeline by default
			HandleOnPipelineSelect(pipelineList, 0, frame)
		})
	}
	// Initial load
	go loadPipelines(bitbucket.BuildQuery(""), false)

	// callback for refresh while watching changes..
	state.PipelineUIState.PipelineList.SetOnRefresh(func() {
		go loadPipelines(bitbucket.BuildQuery(""), false)
	})

	// Animate status with throttled refresh
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS

		defer ticker.Stop()

		for range ticker.C {
			if len(pipelineList) == 0 {
				continue
			}
			frame++

			state.PipelineUIState.App.QueueUpdateDraw(func() {
				for _, pp := range pipelineList {
					status := pp.State.Result.Name
					if status == "" {
						status = pp.State.Name
					}

					if !status.NeedsTracking() {
						continue
					}

					// Check cache
					cached, found := pipelineCache[pp.UUID]
					if found && time.Since(cached.lastFetched) < 10*time.Second {
						// Use cached result
						updated := cached.data
						icon := util.GetIconForStatusWithColorAnimated(updated.State.Name, frame)
						color := util.GetColorForStatus(updated.State.Name)
						text := fmt.Sprintf("%s %s", icon, updated.State.Name)

						state.PipelineUIState.PipelineList.UpdateStatus(updated.UUID, util.CellFormat(text, color))
						continue
					}

					// Fetch updated pipeline status
					updated := bitbucket.FetchPipeline(pp.UUID)
					pipelineCache[pp.UUID] = pipelineCacheEntry{
						lastFetched: time.Now(),
						data:        *updated,
					}

					icon := util.GetIconForStatusWithColorAnimated(updated.State.Name, frame)
					color := util.GetColorForStatus(updated.State.Name)
					text := fmt.Sprintf("%s %s", icon, updated.State.Name)
					state.PipelineUIState.PipelineList.UpdateStatus(updated.UUID, util.CellFormat(text, color))
				}
			})
		}
	}()

	// Infinite scroll - load next page if user scrolls to the bottom
	state.PipelineUIState.PipelineList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		row, _ := state.PipelineUIState.PipelineList.GetSelection()
		totalRows := state.PipelineUIState.PipelineList.GetRowCount()

		if key == tcell.KeyDown && row >= totalRows-2 && nextPageURL != "" {
			log.Printf("[INFO] Near bottom, loading more: %s", nextPageURL)

			query := util.ExtractQueryFromNextURL(nextPageURL)
			nextPageURL = "" // prevent duplicate triggers

			go loadPipelines(query, true)
		}

		return event
	})
}

func HandleOnPipelineSelect(pipelines []types.PipelineResponse, row int, frame int) {
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

	state.SetSelectedPipeline(&selectedPipeline)
	state.PipelineUIState.PipelineList.UpdateSelectedRow(row)

	log.Printf("Selected pipeline UUID: %s, Name: %s", selectedPipeline.UUID, selectedPipeline.Creator.DisplayName)

	support.ShowPipelineLoadingSpinner(state.PipelineUIState.PipelineSteps, func() (interface{}, error) {
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
			support.UpdateView(state.PipelineUIState.PipelineSteps, fmt.Sprintf("[red]Error: %v[-]", err))
			return
		}

		if state.PipelineUIState.TrackingCancelFunc != nil {
			state.PipelineUIState.TrackingCancelFunc()
		}

		// Create new cancellable context for tracker
		ctx, cancel := context.WithCancel(context.Background())
		state.PipelineUIState.TrackingCancelFunc = cancel

		// Track pipeline based on all status change..
		go TrackPipelineLive(ctx, selectedPipeline)

		stepsView := GenerateStepsView(steps, selectedPipeline, frame)

		support.UpdateView(state.PipelineUIState.PipelineStepsDebugView, GeneratePPDebugInfo(selectedPipeline))
		support.UpdateView(state.PipelineUIState.PipelineSteps, stepsView)

		//	HandleOnStepSelect(steps, selectedPipeline, 0) // Auto select first step and fetch the info
	})
}

func TrackPipelineLive(ctx context.Context, pipeline types.PipelineResponse) {
	if !pipeline.State.Name.NeedsTracking() {
		return
	}
	log.Printf("[INFO] Starting live tracking for pipeline: %s", pipeline.UUID)

	animationTicker := time.NewTicker(100 * time.Millisecond) // 10 FPS animation
	defer animationTicker.Stop()

	fetchTicker := time.NewTicker(3 * time.Second) // 3 seconds fetch interval
	defer fetchTicker.Stop()

	const stableThreshold = 3
	stableCount := 0

	frame := 0
	steps := []types.StepDetail{}

	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] Tracking for pipeline %s cancelled", pipeline.UUID)
			return

		case <-animationTicker.C:
			// Just update animation frame, patch UI with same steps but new frame for spinner
			frame++
			state.PipelineUIState.PipelineSteps.PatchSteps(steps, frame)

		case <-fetchTicker.C:
			updatedPipeline := bitbucket.FetchPipeline(pipeline.UUID)
			if updatedPipeline.UUID == "" {
				log.Printf("[WARN] Could not fetch pipeline update for %s", pipeline.UUID)
				continue
			}

			newSteps := bitbucket.FetchPipelineSteps(updatedPipeline.UUID)
			if newSteps == nil {
				log.Printf("[WARN] Could not fetch updated steps for pipeline %s", updatedPipeline.UUID)
				continue
			}

			// Update steps with new data
			steps = newSteps

			// Check if all steps are done
			allDone := true
			for _, step := range steps {
				if step.State.Name.InProgress() || step.State.Name.Running() || step.State.Name.Pending() {
					allDone = false
					break
				}
			}

			if allDone {
				stableCount++
				log.Printf("[INFO] All steps complete for pipeline %s (%d/%d confirmations)", pipeline.UUID, stableCount, stableThreshold)
			} else {
				stableCount = 0
			}

			// Patch UI immediately on fetch as well
			state.PipelineUIState.PipelineSteps.PatchSteps(steps, frame)

			if !updatedPipeline.State.Name.InProgress() && stableCount >= stableThreshold {
				log.Printf("[INFO] Pipeline %s finished and stable. Stopping tracker.", pipeline.UUID)
				return
			}
		}
	}
}

func EmptyAllPipelineListDependentViews() {
	// Always start a fresh slate for rest of dependent views
	emptyView := tview.NewFlex()
	support.UpdateView(state.PipelineUIState.PipelineSteps, emptyView)
	support.UpdateView(state.PipelineUIState.PipelineStepsDebugView, emptyView)
	support.UpdateView(state.PipelineUIState.PipelineStep, emptyView)
	support.UpdateView(state.PipelineUIState.PipelineStepCommandsView, emptyView)
	support.UpdateView(state.PipelineUIState.PipelineStepCommandLogView, emptyView)
}
