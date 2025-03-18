package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AIService interface {
	Chat(prompt string) (*string, error)
}

// AnalyzeProcess calls the AI service to get an explanation of the process's purpose
// and an assessment of whether it might be malicious.
func AnalyzeProcess(pd ProcessData) (AnalysisResult, error) {
	// Build the prompt using process details.
	prompt := fmt.Sprintf(
		"Analyze the following process details and provide a brief explanation of its purpose. Then, determine if the process might be malicious. Respond in JSON format with two keys: 'description' (a brief explanation) and 'malicious' (true or false). Process details: PID: %d, Name: %s, Status: %s, CPU: %.2f%%, Memory: %.2f%%, Command: %s.",
		pd.PID, pd.Name, pd.Status, pd.CPU, pd.Memory, pd.Command,
	)

	var service AIService = Groq{}
	message, err := service.Chat(prompt)
	if err != nil{
		return AnalysisResult{}, err
	}

	*message = strings.TrimSpace(*message)             // Remove leading/trailing whitespace
	*message = strings.TrimPrefix(*message, "```json") // Remove ```json
	*message = strings.TrimSuffix(*message, "```")     // Remove ```

	var analysis AnalysisResult
	err = json.Unmarshal([]byte(*message), &analysis)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return AnalysisResult{}, err
	}

	return analysis, err
}
