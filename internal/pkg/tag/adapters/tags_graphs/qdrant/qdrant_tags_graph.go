package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/usecase/dependencies"
)

var _ dependencies.TagsGraph = &QdrantTagsGraph{}

type Inferer interface {
	Infer(content string) ([]float32, error)
}

const pointsUrl = "http://localhost:6333/collections/tags/points"

type QdrantTagsGraph struct {
	inferer Inferer
}

func NewQdrantTagsGraph(
	inferer Inferer,
) *QdrantTagsGraph {
	return &QdrantTagsGraph{
		inferer: inferer,
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
		pointsUrl,
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
	// TODO
}

// TODO: return ids instead of models.Tag
func (g *QdrantTagsGraph) ListClosestTags(tag *models.Tag, limit uint32) ([]*models.Tag, error) {
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
		fmt.Sprintf("%s/search", pointsUrl),
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

	// TODO: map response

	return []*models.Tag{}, nil
}
