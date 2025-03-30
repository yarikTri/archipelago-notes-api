package ollama

import (
	"fmt"
	"strings"

	"github.com/yarikTri/archipelago-notes-api/internal/clients/llm"
)

type TagSuggesterPort interface {
	SuggestTags(text string) ([]string, error)
}

type TagSuggester struct {
	openAiClient *llm.OpenAiClient
}

func NewTagSuggester(openAiClient *llm.OpenAiClient) *TagSuggester {
	return &TagSuggester{
		openAiClient: openAiClient,
	}
}

func (s *TagSuggester) SuggestTags(text string) ([]string, error) {
	prompt := fmt.Sprintf(`Given the following text, suggest relevant tags (max 5). Return only the tags separated by commas, no other text:
%s`, text)

	response, err := s.openAiClient.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tags: %w", err)
	}

	// Clean up the response and split into tags
	tags := strings.Split(strings.TrimSpace(response), ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	return tags, nil
}
