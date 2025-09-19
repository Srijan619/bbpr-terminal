package pipeline

import (
	"encoding/json"
	"log"
	"simple-git-terminal/types"
)

func GeneratePPDetail(pipeline *types.PipelineResponse) string {
	jsonBytes, err := json.MarshalIndent(pipeline, "", "  ")
	if err != nil {
		log.Printf("Error marshalling pipeline: %v", err)
		return "Error generating pipeline detail"
	}
	return string(jsonBytes)
}
