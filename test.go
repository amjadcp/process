package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

// AnalysisResult represents the analysis outcome from an AI service.
type AnalysisResult struct {
	Description string `json:"description"`
	Malicious   bool   `json:"malicious"`
}

// AIService defines the interface for third-party AI services.
type AIService interface {
	Analyze(pd ProcessData) (AnalysisResult, error)
}

// GetAIService is a factory function that returns an AIService implementation
// based on the provided service name.
func GetAIService(serviceName string) (AIService, error) {
	switch serviceName {
	case "groq":
		return NewGroqService(), nil
	case "openapi":
		return NewOpenAPIService(), nil
	case "ollama":
		return NewOllamaService(), nil
	default:
		return nil, fmt.Errorf("unsupported AI service: %s", serviceName)
	}
}

// AnalyzeProcess is a convenience function that selects the AI service based on the
// environment variable "AI_SERVICE" (defaulting to "groq" if not set) and performs analysis.
func AnalyzeProcess(pd ProcessData) (AnalysisResult, error) {
	serviceName := os.Getenv("AI_SERVICE")
	if serviceName == "" {
		serviceName = "groq"
	}
	service, err := GetAIService(serviceName)
	if err != nil {
		return AnalysisResult{}, err
	}
	return service.Analyze(pd)
}

// ------------------ Groq Service Implementation ------------------

// GroqService implements the AIService interface using the Groq API.
type GroqService struct {
	apiURL string
	apiKey string
}

// NewGroqService returns a new instance of GroqService.
func NewGroqService() *GroqService {
	return &GroqService{
		apiURL: "https://api.groq.com/openai/v1/chat/completions",
		apiKey: os.Getenv("GROQ_API_KEY"),
	}
}

// Analyze sends the process details to the Groq API and returns the analysis result.
func (g *GroqService) Analyze(pd ProcessData) (AnalysisResult, error) {
	if g.apiKey == "" {
		return AnalysisResult{}, fmt.Errorf("GROQ_API_KEY not set")
	}
	prompt := fmt.Sprintf(
		"Analyze the following process details and provide a brief explanation of its purpose. Then, determine if the process might be malicious. Respond in JSON format with two keys: 'description' (a brief explanation) and 'malicious' (true or false). Process details: PID: %d, Name: %s, Status: %s, CPU: %.2f%%, Memory: %.2f%%, Command: %s.",
		pd.PID, pd.Name, pd.Status, pd.CPU, pd.Memory, pd.Command,
	)

	// Build the request payload for the Groq API.
	reqBody := struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}{
		Model: "llama-3.3-70b-versatile",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return AnalysisResult{}, err
	}

	req, err := http.NewRequest("POST", g.apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return AnalysisResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return AnalysisResult{}, err
	}
	defer resp.Body.Close()

	var groqResp struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return AnalysisResult{}, err
	}
	if len(groqResp.Choices) == 0 {
		return AnalysisResult{}, fmt.Errorf("no choices in groq response")
	}

	// Attempt to parse the response content as JSON.
	var analysis AnalysisResult
	err = json.Unmarshal([]byte(groqResp.Choices[0].Message.Content), &analysis)
	if err != nil {
		// If parsing fails, return the raw content as the description and assume it's safe.
		return AnalysisResult{
			Description: groqResp.Choices[0].Message.Content,
			Malicious:   false,
		}, nil
	}
	return analysis, nil
}

// Message represents a single message in the API request payload.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ------------------ Stub Implementations for Other Services ------------------

// NewOpenAPIService returns a stub implementation for an OpenAPI-based AI service.
func NewOpenAPIService() AIService {
	return &StubService{name: "openapi"}
}

// NewOllamaService returns a stub implementation for an Ollama-based AI service.
func NewOllamaService() AIService {
	return &StubService{name: "ollama"}
}

// StubService is a placeholder implementation of AIService for services not yet implemented.
type StubService struct {
	name string
}

// Analyze for StubService returns an error indicating the service is not implemented.
func (s *StubService) Analyze(pd ProcessData) (AnalysisResult, error) {
	return AnalysisResult{}, fmt.Errorf("AI service '%s' not implemented yet", s.name)
}
