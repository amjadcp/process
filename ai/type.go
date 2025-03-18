package ai

// ProcessData contains the process details to be analyzed.
type ProcessData struct {
	PID     int32   `json:"pid"`
	Name    string  `json:"name"`
	Status  string  `json:"status"`
	CPU     float64 `json:"cpu"`
	Memory  float32 `json:"memory"`
	Command string  `json:"command"`
}

// AnalysisResult represents the analysis outcome from the AI service.
type AnalysisResult struct {
	Description string `json:"description"`
	Malicious   bool   `json:"malicious"`
}

// AIService defines the interface for AI interactions.
type AIService interface {
	Chat(prompt string) (string, error)
}