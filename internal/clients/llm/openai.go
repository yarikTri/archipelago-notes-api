package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAiClient struct {
	baseURL string
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func NewOpenAiClient(baseURL string) *OpenAiClient {
	return &OpenAiClient{
		baseURL: baseURL,
	}
}

func (c *OpenAiClient) Generate(model string, prompt string, stream bool, system string) (string, error) {
	reqBody := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: stream,
		System: system,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(c.baseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var generateResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return generateResp.Response, nil
}
