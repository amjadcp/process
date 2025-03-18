package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amjadcp/process/config"
)

type Groq struct {
	URL string
	Model string
	APIKEY string
}


func (g Groq)Chat(prompt string)(*string, error) {
	reqBody := GroqRequest{
		Model: config.Env.GROQ_MODEL,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", config.Env.GROQ_API_URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Env.GROQ_API_KEY)

	client := http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, err
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in groq response")
	}

	return &groqResp.Choices[0].Message.Content, nil
}



