package triton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TritonClient struct {
	baseURL string
}

type InferenceRequest struct {
	Inputs []Input `json:"inputs"`
}

type Input struct {
	Name     string        `json:"name"`
	Datatype string        `json:"datatype"`
	Shape    []int         `json:"shape"`
	Data     []interface{} `json:"data"`
}

type InferenceResponse struct {
	Outputs []Output `json:"outputs"`
}

type Output struct {
	Name     string    `json:"name"`
	Datatype string    `json:"datatype"`
	Shape    []int     `json:"shape"`
	Data     []float32 `json:"data"`
}

func NewTritonClient(baseURL string) *TritonClient {
	return &TritonClient{
		baseURL: baseURL,
	}
}

func (c *TritonClient) GetEmbedding(text string) ([]float32, error) {
	reqBody := InferenceRequest{
		Inputs: []Input{
			{
				Name:     "text_feature",
				Datatype: "BYTES",
				Shape:    []int{1, 1},
				Data:     []interface{}{text},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v2/models/ensemble_model/infer", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var inferenceResp InferenceResponse
	if err := json.Unmarshal(body, &inferenceResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(inferenceResp.Outputs) == 0 {
		return nil, fmt.Errorf("no outputs in response")
	}

	return inferenceResp.Outputs[0].Data, nil
}
