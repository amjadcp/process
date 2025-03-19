package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/amjadcp/process-monitor/config"
)

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

func NewAIService(aiService string) AIService {
	switch aiService {
	case "groq":
		return &Groq{
			URL:    config.Env.GROQ_API_URL,
			Model:  config.Env.GROQ_MODEL,
			APIKEY: config.Env.GROQ_API_KEY,
		}
	case "ollama":
		return &Ollama{
			URL:    config.Env.OLLAMA_API_URL,
			Model:  config.Env.OLLAMA_MODEL,
			APIKEY: config.Env.OLLAMA_API_KEY,
		}
	default:
		return nil
	}
}

// AnalyzeProcess calls the AI service to explain the process’s purpose
// and to assess whether it might be malicious.
func AnalyzeProcess(pd ProcessData) (AnalysisResult, error) {
	prompt := fmt.Sprintf(
		"Analyze the following process details and provide a brief explanation of its purpose. Then, determine if the process might be malicious. Respond in JSON format with two keys: 'description' (a brief explanation) and 'malicious' (true or false). Process details: PID: %d, Name: %s, Status: %s, CPU: %.2f%%, Memory: %.2f%%, Command: %s.",
		pd.PID, pd.Name, pd.Status, pd.CPU, pd.Memory, pd.Command,
	)

	// Inject configuration using dependency injection.
	// var service AIService = &Groq{
	// 	URL:    config.Env.GROQ_API_URL,
	// 	Model:  config.Env.GROQ_MODEL,
	// 	APIKEY: config.Env.GROQ_API_KEY,
	// }
	service := NewAIService(config.Env.AI_SERVICE) 
	message, err := service.Chat(prompt)
	if err != nil {
		return AnalysisResult{}, err
	}

	message = strings.TrimSpace(message)
	message = strings.TrimPrefix(message, "```json")
	message = strings.TrimSuffix(message, "```")

	var analysis AnalysisResult
	err = json.Unmarshal([]byte(message), &analysis)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return AnalysisResult{}, err
	}

	return analysis, nil
}
