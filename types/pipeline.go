package types

type BitbucketPipelineResponse struct {
	Values     []PipelineResponse `json:"values"`
	Pagination Pagination         `json:"pagination"`
}

type PipelineResponse struct {
	UUID             string            `json:"uuid"`
	State            State             `json:"state"`
	Duration         int               `json:"duration_in_seconds"`
	CreatedOn        string            `json:"created_on"`
	CompletedOn      string            `json:"completed_on"`
	BuildNumber      int               `json:"build_number"`
	BuildSecondsUsed int               `json:"build_seconds_used"`
	FirstSuccessful  bool              `json:"first_successful"`
	RunNumber        int               `json:"run_number"`
	Creator          User              `json:"creator"`
	Target           PipelineRefTarget `json:"target"`
	Trigger          Trigger           `json:"trigger"`
}

type PipelineRefTarget struct {
	RefType  string   `json:"ref_type"`
	RefName  string   `json:"ref_name"`
	Selector Selector `json:"selector"`
	Commit   Commit   `json:"commit"`
	Type     string   `json:"type"`
}

type Trigger struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

type Service struct {
	Name string `json:"name"`
}

type Selector struct {
	Type string `json:"type"`
}

type CommitLinks struct {
	Self Link `json:"self"`
	HTML Link `json:"html"`
}

type Link struct {
	Href string `json:"href"`
}

type State struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Result Result `json:"result"`
}

type Result struct {
	Name PipelineStatus `json:"name"`
	Type string         `json:"type"`
}

// Pipeline Status

type PipelineStatus string

const (
	StatusParsing  PipelineStatus = "PARSING"
	StatusPending  PipelineStatus = "PENDING"
	StatusPaused   PipelineStatus = "PAUSED"
	StatusHalted   PipelineStatus = "HALTED"
	StatusBuilding PipelineStatus = "BUILDING"
	StatusError    PipelineStatus = "ERROR"
	StatusPassed   PipelineStatus = "PASSED"
	StatusFailed   PipelineStatus = "FAILED"
	StatusStopped  PipelineStatus = "STOPPED"
	StatusUnknown  PipelineStatus = "UNKNOWN"
	InProgress     PipelineStatus = "IN_PROGRESS"
	Successful     PipelineStatus = "SUCCESSFUL"
)

// Attach methods to this type (must be in same package!)
func (s PipelineStatus) Failed() bool     { return s == StatusFailed }
func (s PipelineStatus) Passed() bool     { return s == StatusPassed }
func (s PipelineStatus) Pending() bool    { return s == StatusPending }
func (s PipelineStatus) Error() bool      { return s == StatusError }
func (s PipelineStatus) Running() bool    { return s == StatusBuilding }
func (s PipelineStatus) Stopped() bool    { return s == StatusStopped }
func (s PipelineStatus) InProgress() bool { return s == InProgress }
func (s PipelineStatus) Successful() bool { return s == Successful }
func (s PipelineStatus) Unknown() bool    { return s == StatusUnknown }
