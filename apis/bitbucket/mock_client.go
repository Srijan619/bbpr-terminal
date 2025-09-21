package bitbucket

import (
	"log"
	"simple-git-terminal/types"
	"sync"
	"time"
)

// Global variables to track mock states and delays
var (
	mockPipelinePollCount    = 0
	pipelinePollTimestamps   = make(map[string]time.Time)
	pipelineStepProgressions = map[string][]StepMock{
		"mock-uuid-1": {
			{
				UUID: "step-1",
				Name: "Install Dependencies",
				States: []types.State{
					{Name: types.StatusPending},
					{Name: types.InProgress},
					{Name: types.StatusPassed, Result: types.Result{Name: types.Successful}},
				},
			},
			{
				UUID: "step-2",
				Name: "Build",
				States: []types.State{
					{Name: types.StatusPending},
					{Name: types.StatusPending},
					{Name: types.InProgress},
					{Name: types.StatusPassed, Result: types.Result{Name: types.Successful}},
				},
			},
		},
		"mock-uuid-2": {
			{
				UUID: "step-1",
				Name: "Run Tests",
				States: []types.State{
					{Name: types.StatusPending},
					{Name: types.InProgress},
					{Name: types.StatusFailed, Result: types.Result{Name: types.StatusFailed}},
				},
			},
			{
				UUID: "step-2",
				Name: "Deploy",
				States: []types.State{
					{Name: types.StatusPending},
					{Name: types.StatusPending},
					{Name: types.StatusPending},
					{Name: types.InProgress},
					{Name: types.StatusPassed, Result: types.Result{Name: types.Successful}},
				},
			},
		},
	}
	pipelineMutex = sync.Mutex{}
)

func TestFetchPipelinesByQuery(query string) ([]types.PipelineResponse, types.Pagination) {
	pipelineMutex.Lock()
	defer pipelineMutex.Unlock()

	mockPipelinePollCount++
	now := time.Now()

	// Update timestamps if not set
	for uuid := range pipelineStepProgressions {
		if _, exists := pipelinePollTimestamps[uuid]; !exists {
			pipelinePollTimestamps[uuid] = now
		}
	}

	log.Printf("[MOCK CLIENT] Fetching Pipelines (poll #%d) with query: %v", mockPipelinePollCount, query)

	// Calculate pollCount per pipeline based on elapsed time for accurate states
	var pipelines []types.PipelineResponse

	for uuid, steps := range pipelineStepProgressions {
		lastPoll := pipelinePollTimestamps[uuid]
		elapsed := now.Sub(lastPoll)
		pollCount := int(elapsed / pollInterval)

		overallState, overallResult := derivePipelineState(steps, pollCount)

		refName := "main"
		if uuid == "mock-uuid-2" {
			refName = "feature/login"
		}

		pipelines = append(pipelines, types.PipelineResponse{
			UUID: uuid,
			State: types.State{
				Name:   overallState,
				Result: types.Result{Name: overallResult},
			},
			Target:    types.PipelineRefTarget{RefName: refName},
			CreatedOn: "2025-09-20T12:00:00Z",
		})
	}

	// Pagination is fixed for mocks, can be expanded if needed
	pagination := types.Pagination{
		Size:    len(pipelines),
		PageLen: 10,
		Page:    1,
		Next:    "",
	}

	return pipelines, pagination
}

func derivePipelineState(steps []StepMock, pollCount int) (types.PipelineStatus, types.PipelineStatus) {
	allPassed := true
	anyFailed := false
	anyInProgress := false

	for _, step := range steps {
		stateIndex := pollCount
		if stateIndex >= len(step.States) {
			stateIndex = len(step.States) - 1
		}
		s := step.States[stateIndex]

		switch s.Name {
		case types.StatusFailed:
			anyFailed = true
			allPassed = false
		case types.StatusPassed:
			// no-op; still might allPassed = true
		case types.InProgress:
			anyInProgress = true
			allPassed = false
		default:
			allPassed = false
		}

		if s.Result.Name == types.StatusFailed {
			anyFailed = true
		}
	}

	switch {
	case anyFailed:
		return types.StatusPassed, types.StatusFailed
	case anyInProgress:
		return types.InProgress, ""
	case allPassed:
		return types.StatusPassed, types.Successful
	default:
		return types.StatusPending, ""
	}
}

// StepMock defines the structure for each step in a pipeline
type StepMock struct {
	UUID   string
	Name   string
	States []types.State
}

// SimulatedFetchPipelineSteps simulates fetching pipeline steps with artificial delays
func SimulatedFetchPipelineSteps(pipelineUUID string) []types.StepDetail {
	pipelineMutex.Lock()
	defer pipelineMutex.Unlock()

	now := time.Now()
	lastPoll, exists := pipelinePollTimestamps[pipelineUUID]
	if !exists {
		pipelinePollTimestamps[pipelineUUID] = now
		lastPoll = now
	}

	elapsed := now.Sub(lastPoll)
	pollCount := int(elapsed / pollInterval)

	log.Printf("[MOCK SIM] Fetching steps for pipeline %s at pollCount %d", pipelineUUID, pollCount)

	steps, found := pipelineStepProgressions[pipelineUUID]
	if !found {
		return nil
	}

	var stepDetails []types.StepDetail
	for _, step := range steps {
		stateIndex := pollCount
		if stateIndex >= len(step.States) {
			stateIndex = len(step.States) - 1
		}
		stepDetails = append(stepDetails, types.StepDetail{
			UUID:  step.UUID,
			Name:  step.Name,
			State: step.States[stateIndex],
		})
	}

	return stepDetails
}

// TestFetchPipeline simulates fetching a single pipeline with throttling and progressive state
func TestFetchPipeline(pipelineUUID string) *types.PipelineResponse {
	pipelineMutex.Lock()
	defer pipelineMutex.Unlock()

	now := time.Now()
	lastPoll, exists := pipelinePollTimestamps[pipelineUUID]
	if !exists {
		pipelinePollTimestamps[pipelineUUID] = now
		lastPoll = now
	}

	elapsed := now.Sub(lastPoll)
	pollCount := int(elapsed / pollInterval)

	log.Printf("[MOCK CLIENT] Fetching single pipeline '%s' at pollCount %d", pipelineUUID, pollCount)

	steps, found := pipelineStepProgressions[pipelineUUID]
	if !found {
		return nil
	}

	overallState := types.StatusPending
	overallResult := types.PipelineStatus("")

	allPassed := true
	anyFailed := false
	anyInProgress := false

	for _, step := range steps {
		stateIndex := pollCount
		if stateIndex >= len(step.States) {
			stateIndex = len(step.States) - 1
		}
		s := step.States[stateIndex]

		switch s.Name {
		case types.StatusFailed:
			anyFailed = true
			allPassed = false
		case types.StatusPassed:
			// Keep allPassed true only if all steps passed
		case types.InProgress:
			anyInProgress = true
			allPassed = false
		default:
			allPassed = false
		}

		if s.Result.Name == types.StatusFailed {
			anyFailed = true
		}
	}

	switch {
	case anyFailed:
		overallState = types.StatusPassed
		overallResult = types.StatusFailed
	case anyInProgress:
		overallState = types.InProgress
	case allPassed:
		overallState = types.StatusPassed
		overallResult = types.Successful
	default:
		overallState = types.StatusPending
	}

	return &types.PipelineResponse{
		UUID: pipelineUUID,
		State: types.State{
			Name:   overallState,
			Result: types.Result{Name: overallResult},
		},
		Target:    types.PipelineRefTarget{RefName: "main"},
		CreatedOn: "2025-09-20T12:00:00Z",
	}
}

// ResetMockState resets all mock states and timestamps
func ResetMockState() {
	pipelineMutex.Lock()
	defer pipelineMutex.Unlock()

	mockPipelinePollCount = 0
	pipelinePollTimestamps = make(map[string]time.Time)
}

// pollInterval defines the artificial delay between each poll
const pollInterval = 5 * time.Second
