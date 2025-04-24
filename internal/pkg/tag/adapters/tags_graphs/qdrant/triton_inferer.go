package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var _ Inferer = &TritonInferer{}

type TritonInferer struct{}

func NewTritonInferer() *TritonInferer {
	return &TritonInferer{}
}

type inferRequest struct {
	Inputs []inferRequestInput `json:"inputs"`
}

type inferRequestInput struct {
	Name     string   `json:"name"`
	Datatype string   `json:"datatype"`
	Shape    []int    `json:"shape"`
	Data     []string `json:"data"`
}

type inferResponse struct {
	Vector []float32 `json:"vector"`
}

func (i *TritonInferer) Infer(content string) ([]float32, error) {
	req := inferRequest{
		Inputs: []inferRequestInput{
			{
				Name:     "text_feature",
				Datatype: "BYTES",
				Shape:    []int{1, 1},
				Data:     []string{content},
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := http.Post(
		// TODO: change port
		"http://localhost:1234/v2/models/ensemble_model/infer",
		"application/json",
		bytes.NewReader(reqJson),
	)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("got response status code %d", httpResp.StatusCode)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp inferResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Vector, nil
}
