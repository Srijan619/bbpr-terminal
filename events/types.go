package events

import "simple-git-terminal/types"

// PipelineSelectedEvent is emitted when the user selects a pipeline
type PipelineSelectedEvent struct {
	Pipeline types.PipelineResponse
	Row      int
}

// StepsUpdatedEvent is emitted when the steps for a pipeline have been fetched
type StepsUpdatedEvent struct {
	PipelineUUID string
	Steps        []types.StepDetail
}
