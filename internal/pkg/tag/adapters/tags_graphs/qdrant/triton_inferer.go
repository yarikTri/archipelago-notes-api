package qdrant

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var _ Inferer = &TritonInferer{}

type TritonInferer struct {
	inferUrl string
}

func NewTritonInferer(host, port string) *TritonInferer {
	if host == "" {
		panic("TritonInferer: host cant be empty")
	}
	if _, err := strconv.Atoi(port); err != nil {
		panic(fmt.Sprintf("TritonInferer: invalid port: %v", err))
	}
	return &TritonInferer{
		inferUrl: fmt.Sprintf("http://%s:%s/v2/models/ensemble_model/infer", host, port),
	}
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

type inferenceResponse struct {
	ModelName    string                    `json:"model_name"`
	ModelVersion string                    `json:"model_version"`
	Parameters   map[string]interface{}    `json:"parameters"`
	Outputs      []inferenceResponseOutput `json:"outputs"`
}

var ErrVectorInResponseNotFound = errors.New("response contains zero outputs")

func (r *inferenceResponse) getFirstVector() ([]float32, error) {

	if len(r.Outputs) == 0 {
		return nil, ErrVectorInResponseNotFound
	}

	return r.Outputs[0].Data, nil
}

// Can be replaced for Parameters in InferenceResponse
// .
// type Parameters struct {
//     SequenceID    int  `json:"sequence_id"`
//     SequenceStart bool `json:"sequence_start"`
//     SequenceEnd   bool `json:"sequence_end"`
// }

type inferenceResponseOutput struct {
	Name     string    `json:"name"`
	Datatype string    `json:"datatype"`
	Shape    []int     `json:"shape"`
	Data     []float32 `json:"data"`
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

	// TODO: make port and host custom.
	httpResp, err := http.Post(
		i.inferUrl,
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

	var resp inferenceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.getFirstVector()
}
