package ollama

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/yarikTri/archipelago-notes-api/internal/clients/llm"
)

const (
	TagSytemPrompt = "Придумай тег, который наиболее точно описывает текст этой заметки - одно слово"
	MaxAttempts    = 2
)

type TagSuggesterPort interface {
	SuggestTags(text string, tagsNum *int) ([]string, error)
}

type OpenAIClientError struct {
	err error
}

func (e *OpenAIClientError) Error() string {
	return fmt.Sprintf("OpenAI client error: %v", e.err)
}

func (e *OpenAIClientError) Unwrap() error {
	return e.err
}

type TagSuggester struct {
	openAiClient          *llm.OpenAiClient
	defaultGenerateTagNum int
	model                 string
}

func NewTagSuggester(openAiClient *llm.OpenAiClient, defaultGenerateTagNum int, model string) *TagSuggester {
	return &TagSuggester{
		openAiClient:          openAiClient,
		defaultGenerateTagNum: defaultGenerateTagNum,
		model:                 model,
	}
}

func isValidTag(tag string) bool {
	// Check for white spaces.
	parts := strings.Fields(tag)
	if len(parts) != 1 {
		return false
	}
	tag = parts[0]

	// Empty or whitespace-only tags
	if len(strings.TrimSpace(tag)) == 0 {
		return false
	}

	// invalidResponses := []string{
	// 	// Linking words
	// 	"это",
	// 	"вот",
	// 	"или",
	// 	"либо",
	// 	"также",
	// 	"еще",
	// 	"затем",
	// 	"далее",
	// 	"итак",
	// 	"значит",
	// }

	// tagLower := strings.ToLower(tag)
	// for _, invalid := range invalidResponses {
	// 	if strings.Contains(tagLower, invalid) {
	// 		return false
	// 	}
	// }

	if len(tag) > 80 || len(tag) < 2 {
		return false
	}

	return true
}

func cleanupTag(response string) string {
	tag := strings.TrimSpace(response)

	// Compile the regular expression
	// This matches any character that is NOT a letter (a-z, A-Z) or number (0-9)
	reg := regexp.MustCompile(`[^a-zA-Z0-9]`)

	// Replace all matched characters with an empty string
	tag = reg.ReplaceAllString(tag, "")

	return tag
}

func (s *TagSuggester) generateOneTagWithRetry(text string) (string, error) {
	for attempt := 0; attempt < MaxAttempts; attempt++ {
		response, err := s.openAiClient.Generate(
			s.model,
			text,
			false, // stream
			TagSytemPrompt,
		)
		if err != nil {
			return "", &OpenAIClientError{err: err}
		}

		tag := cleanupTag(response)

		if isValidTag(tag) {
			fmt.Printf("Got response from ollama (approved): %s\n", response)
			return strings.ToLower(tag), nil
		}
		fmt.Printf("Got response from ollama (not approved): %s\n", response)
	}
	return "", fmt.Errorf("failed to generate valid tag after %d attempts", MaxAttempts)
}

func (s *TagSuggester) SuggestTags(text string, tagsNum *int) ([]string, error) {
	numTags := s.defaultGenerateTagNum
	if tagsNum != nil {
		numTags = *tagsNum
	}

	tags := make([]string, 0, numTags)
	for i := 0; i < numTags; i++ {
		tag, err := s.generateOneTagWithRetry(text)
		if err != nil {
			var clientErr *OpenAIClientError
			if errors.As(err, &clientErr) {
				return nil, err // Return immediately on client errors
			}
			if len(tags) > 0 {
				return tags, nil // Return partial results on validation errors
			}
			return nil, fmt.Errorf("failed to generate tag %d: %w", i+1, err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
