package ollama

import (
	"fmt"
	"strings"

	"github.com/yarikTri/archipelago-notes-api/internal/clients/llm"
)

const (
	TagSytemPrompt = "Придумай тег, который наиболее точно описывает текст этой заметки - одно слово"
)

type TagSuggesterPort interface {
	SuggestTags(text string, tagsNum *int) ([]string, error)
}

type TagSuggester struct {
	openAiClient          *llm.OpenAiClient
	defaultGenerateTagNum int
	system                string
}

func NewTagSuggester(openAiClient *llm.OpenAiClient, defaultGenerateTagNum int, system string) *TagSuggester {
	return &TagSuggester{
		openAiClient:          openAiClient,
		defaultGenerateTagNum: defaultGenerateTagNum,
		system:                system,
	}
}

func (s *TagSuggester) suggestTag(text string) (string, error) {
	response, err := s.openAiClient.Generate("llama2", text, false, s.system)
	if err != nil {
		return "", fmt.Errorf("failed to generate tag: %w", err)
	}

	// Clean up the response and get first tag
	tag := strings.TrimSpace(strings.Split(response, ",")[0])
	return tag, nil
}

func (s *TagSuggester) SuggestTags(text string, tagsNum *int) ([]string, error) {
	numTags := s.defaultGenerateTagNum
	if tagsNum != nil {
		numTags = *tagsNum
	}

	tags := make([]string, 0, numTags)
	for i := 0; i < numTags; i++ {
		tag, err := s.suggestTag(text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tag %d: %w", i+1, err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
