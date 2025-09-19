package types

type BitbucketStepsResponse struct {
	Page    int          `json:"page"`
	Values  []StepDetail `json:"values"`
	Size    int          `json:"size"`
	PageLen int          `json:"pagelen"`
}

type StepDetail struct {
	UUID                    string               `json:"uuid"`
	Name                    string               `json:"name"`
	Trigger                 Trigger              `json:"trigger"`
	State                   State                `json:"state"`
	StartedOn               string               `json:"started_on"`
	CompletedOn             string               `json:"completed_on"`
	DurationInSeconds       int                  `json:"duration_in_seconds"`
	RunNumber               int                  `json:"run_number"`
	BuildSecondsUsed        int                  `json:"build_seconds_used"`
	TestReportDefinition    TestReportDefinition `json:"testReportDefinition"`
	Tasks                   Tasks                `json:"tasks"`
	Arch                    string               `json:"arch"`
	MaxTime                 int                  `json:"maxTime"`
	IsArtifactsDownloadable bool                 `json:"is_artifacts_download_enabled"`
	Type                    string               `json:"type"`

	// Optional fields you might encounter, include only if needed
	// Services             []Service      `json:"services"`
	// Caches               []Cache        `json:"caches"`
	// ResourceLimits       ResourceLimits `json:"resource_limits"`
	// TestReport           TestReport     `json:"test_report"`
	// FeatureFlags         []FeatureFlag  `json:"feature_flags"`
	//
	// 👇 Add these missing fields (they exist in JSON)
	SetupCommands    []CommandDetail `json:"setup_commands"`
	ScriptCommands   []CommandDetail `json:"script_commands"`
	TeardownCommands []CommandDetail `json:"teardown_commands"`
	Pipeline         PipelineInfo    `json:"pipeline"` // nested pipeline object
}

type Cache struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type ResourceLimits struct {
	MemoryLimitInMB  int `json:"memory_limit_in_megabytes"`
	CPULimitInMillis int `json:"cpu_limit_in_millicores"`
}

type TestReportDefinition struct {
	TargetDirectories  []string `json:"target_directories"`
	IgnoredDirectories []string `json:"ignored_directories"`
	SearchDepth        int      `json:"search_depth"`
	Paths              []string `json:"paths"`
	IgnorePaths        []string `json:"ignore_paths"`
	CaptureOn          string   `json:"capture_on"`
}

type TestReport struct {
	Definition TestReportDefinition `json:"definition"`
	Result     TestReportResult     `json:"result"`
}

type TestReportResult struct {
	StepUUID     string `json:"step_uuid"`
	PipelineUUID string `json:"pipeline_uuid"`
	Type         string `json:"type"`
}

type FeatureFlag struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // Can be bool, int, string, etc.
}

type Tasks struct {
	ExecutionPhases ExecutionPhases `json:"execution_phases"`
}

type ExecutionPhases struct {
	Setup    []Phase `json:"SETUP"`
	Main     []Phase `json:"MAIN"`
	Teardown []Phase `json:"TEARDOWN"`
}

type Phase struct {
	Environment []interface{} `json:"environment"` // Can define type if needed
	Commands    []Command     `json:"commands"`
}

type Command struct {
	CommandString string `json:"command_string"`
}

type Image struct {
	Name string `json:"name"`
}

type PipelineInfo struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

type CommandDetail struct {
	CommandType string `json:"commandType"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Action      string `json:"action,omitempty"`
}
