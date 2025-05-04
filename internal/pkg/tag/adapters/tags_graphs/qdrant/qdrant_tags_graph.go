package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/usecase/dependencies"
)

var _ dependencies.TagsGraph = &QdrantTagsGraph{}

type Inferer interface {
	Infer(content string) ([]float32, error)
}

// const pointsUrl = "http://localhost:6333/collections/tags/points"

type QdrantTagsGraph struct {
	pointsUrl string
	inferer   Inferer
}

func NewQdrantTagsGraph(
	inferer Inferer,
	host string,
	port string,
) *QdrantTagsGraph {
	if host == "" {
		panic("QdrantTagsGraph: host cant be empty")
	}
	if _, err := strconv.Atoi(port); err != nil {
		panic(fmt.Sprintf("QdrantTagsGraph: invalid port: %v", err))
	}

	return &QdrantTagsGraph{
		inferer:   inferer,
		pointsUrl: fmt.Sprintf("http://%s:%s/collections/tags/points", host, port),
	}
}

type updateOrCreateTagRequest struct {
	Points []updateOrCreateTagRequestPoint `json:"points"`
}

type updateOrCreateTagRequestPoint struct {
	Id      string                               `json:"id"`
	Payload updateOrCreateTagRequestPointPayload `json:"payload"`
	Vector  []float32                            `json:"vector"`
}

type updateOrCreateTagRequestPointPayload struct {
	Tag    string `json:"tag"`
	UserId string `json:"user_id"`
}

func (g *QdrantTagsGraph) UpdateOrCreateTag(tag *models.Tag) error {
	inferVec, err := g.inferer.Infer(tag.Name)
	if err != nil {
		return fmt.Errorf("failed to infer tag %s: %w", tag.Name, err)
	}

	req := updateOrCreateTagRequest{
		Points: []updateOrCreateTagRequestPoint{
			{
				Id: tag.ID.String(),
				Payload: updateOrCreateTagRequestPointPayload{
					Tag:    tag.Name,
					UserId: tag.UserID.String(),
				},
				Vector: inferVec,
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(
		http.MethodPut,
		g.pointsUrl,
		bytes.NewReader(reqJson),
	)
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return fmt.Errorf("got response status code %d", httpResp.StatusCode)
	}

	return nil
}

type listClosestTagsRequest struct {
	Vector []float32                    `json:"query"`
	Limit  uint32                       `json:"limit"`
	Filter listClosestTagsRequestFilter `json:"filter"`
}

type listClosestTagsRequestFilter struct {
	Must []listClosestTagsRequestFilterMust `json:"must"`
}

type listClosestTagsRequestFilterMust struct {
	Key   string            `json:"key"`
	Match map[string]string `json:"match"`
}

// type listClosestTagsRequestFilterMustMatch struct {
// 	Value string `json:"value"`
// }

type listClosestTagsResponse struct {
	Result []struct {
		ID string `json:"id"`
		// Version   uint64    `json:"version"`
		// Score     float64   `json:"score"`
		// Payload   *Payload  `json:"payload,omitempty"`
		// Vector    []float64 `json:"vector,omitempty"`
	} `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (r *listClosestTagsResponse) error() error {
	if r.Status == "ok" {
		return nil
	}

	return fmt.Errorf("qdrant returned error: %s", r.Status)
}

func (r *listClosestTagsResponse) getAllIds() ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(r.Result))

	for _, v := range r.Result {
		parsedID, err := uuid.FromString(v.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, parsedID)
	}

	return ids, nil
}

func (g *QdrantTagsGraph) ListClosestTagsIds(tag *models.Tag, limit uint32) ([]uuid.UUID, error) {
	inferVec, err := g.inferer.Infer(tag.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to infer tag %s: %w", tag.Name, err)
	}

	req := listClosestTagsRequest{
		Vector: inferVec,
		Limit:  limit,
		Filter: listClosestTagsRequestFilter{
			Must: []listClosestTagsRequestFilterMust{
				{
					Key: tag.ID.String(),
					Match: map[string]string{
						"user_id": tag.UserID.String(),
					},
				},
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := http.Post(
		fmt.Sprintf("%s/search", g.pointsUrl),
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

	var resp listClosestTagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if err := resp.error(); err != nil {
		return nil, err
	}

	return resp.getAllIds()
}

type deleteByIDRequest struct {
	Points []string `json:"points"`
}

func newDeleteByIDRequest(id uuid.UUID) deleteByIDRequest {
	return deleteByIDRequest{
		Points: []string{id.String()},
	}
}

type deleteByIDResponseOperationResult struct {
	Status      string `json:"status"`
	OperationID int64  `json:"operation_id"`
}

type deleteByIDResponse struct {
	// Usage  UsageStats      `json:"usage"`
	Time   float64                           `json:"time"`
	Status string                            `json:"status"`
	Result deleteByIDResponseOperationResult `json:"result"`
}

func (r *deleteByIDResponse) error() error {
	if r.Status == "ok" {
		return nil
	}

	return fmt.Errorf("qdrant returned error: %s", r.Status)
}

const WAIT_FOR_PROCESS = true

func (g *QdrantTagsGraph) DeleteByID(tagID uuid.UUID) error {
	req := newDeleteByIDRequest(tagID)

	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(
		fmt.Sprintf("%s/delete?wait=%t", g.pointsUrl, WAIT_FOR_PROCESS),
		"application/json",
		bytes.NewReader(reqJson),
	)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return fmt.Errorf("got response status code %d", httpResp.StatusCode)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	var resp deleteByIDResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}

	if err := resp.error(); err != nil {
		return err
	}

	return nil
}
