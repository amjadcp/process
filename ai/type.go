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

// GroqRequest represents the payload sent to the Groq API.
type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a single message in the Groq API request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GroqResponse represents the expected structure of the Groq API response.
type GroqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}