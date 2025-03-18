package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Groq implements the AIService interface using the Ollama API.
type Ollama struct {
	URL    string
	Model  string
	APIKEY string
}

// OllamaRequest represents the payload sent to the Ollama API.
type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream bool `json:"stream"`
}

// Message represents a single message in the Ollama API request.
type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaResponse represents the expected structure of the Ollama API response.
type OllamaResponse struct {
	Model             string `json:"model"`
	CreatedAt         string `json:"created_at"`
	Message           struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	DoneReason        string `json:"done_reason"`
	Done              bool   `json:"done"`
	TotalDuration     int64  `json:"total_duration"`
	LoadDuration      int64  `json:"load_duration"`
	PromptEvalCount   int    `json:"prompt_eval_count"`
	PromptEvalDuration int64 `json:"prompt_eval_duration"`
	EvalCount         int    `json:"eval_count"`
	EvalDuration      int64  `json:"eval_duration"`
}


// Chat sends a prompt to the Ollama API and returns the response as a string.
func (o *Ollama) Chat(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model: o.Model,
		Messages: []OllamaMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", o.URL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.APIKEY)

	client := http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", err
	}
	return ollamaResp.Message.Content, nil
}
