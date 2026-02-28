package common

// LabMetadata contains information about a single lab
type LabMetadata struct {
	ID     string
	Name   string
	Charts map[string]ChartMetadata
}

// GetLabsResponse returns all available labs for the UI
type GetLabsResponse struct {
	Labs []LabMetadata `json:"labs"`
}

// LabConfig is the complete configuration for a lab
type LabConfig struct {
	Lab    LabMetadata      `json:"lab"`
	Charts map[string]Chart `json:"charts"`
}
