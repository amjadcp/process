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

// Chat sends a prompt to the Groq API and returns the response as a string.
func (g Groq) Chat(prompt string) (string, error) {
	reqBody := GroqRequest{
		Model: g.Model,
		Messages: []Message{
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
