package events

import "simple-git-terminal/types"

// PipelineSelectedEvent is emitted when the user selects a pipeline
type PipelineSelectedEvent struct {
	Pipeline types.PipelineResponse
	Row      int
}

type StepSelectedEvent struct {
	Step types.StepDetail
	Row  int
}

// StepsUpdatedEvent is emitted when the steps for a pipeline have been fetched
type StepsUpdatedEvent struct {
	PipelineUUID string
	Steps        []types.StepDetail
}
