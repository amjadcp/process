package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Groq implements the AIService interface using the Groq API.
type Groq struct {
	URL    string
	Model  string
	APIKEY string
}

// GroqRequest represents the payload sent to the Groq API.
type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []GroqMessage `json:"messages"`
}

// Message represents a single message in the Groq API request.
type GroqMessage struct {
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


// Chat sends a prompt to the Groq API and returns the response as a string.
func (g *Groq) Chat(prompt string) (string, error) {
	reqBody := GroqRequest{
		Model: g.Model,
		Messages: []GroqMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", g.URL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.APIKEY)

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in groq response")
	}

	return groqResp.Choices[0].Message.Content, nil
}
