package ollama

import (
	"errors"
	"fmt"
	"strconv"
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

// Common validation rules for LLM-generated tags
func isValidTag(tag string) bool {
	// Empty or whitespace-only tags
	if len(strings.TrimSpace(tag)) == 0 {
		return false
	}

	// Single words that cannot be valid tags (common LLM responses and meta-words)
	invalidResponses := []string{
		// LLM self-references and apologies
		"извините",
		"простите",
		"прошу",
		"понимаю",
		"помогу",
		"попробую",
		"постараюсь",
		"проанализирую",
		"рассмотрим",

		// Common meta-responses
		"тег",
		"тэг",
		"метка",
		"категория",
		"раздел",
		"тема",
		"описание",
		"ключевое",
		"основное",
		"главное",

		// Linking words
		"это",
		"вот",
		"или",
		"либо",
		"также",
		"еще",
		"затем",
		"далее",
		"итак",
		"значит",

		// Common LLM fillers
		"пожалуйста",
		"конечно",
		"возможно",
		"вероятно",
		"например",
		"допустим",
		"предположим",
		"следовательно",
		"соответственно",

		// Analysis words
		"анализ",
		"обзор",
		"рассмотрение",
		"изучение",
		"исследование",

		// Task-related
		"задача",
		"цель",
		"задание",
		"требование",
		"результат",

		// Common non-tag responses
		"текст",
		"документ",
		"содержание",
		"информация",
		"данные",
		"материал",
		"контент",
		"заметка",
		"статья",
		"запись",

		// Confirmation words
		"хорошо",
		"ладно",
		"понятно",
		"согласен",
		"подтверждаю",
	}

	tagLower := strings.ToLower(tag)
	for _, invalid := range invalidResponses {
		if strings.Contains(tagLower, invalid) {
			return false
		}
	}

	// Check for reasonable tag length (adjust as needed)
	if len(tag) > 50 || len(tag) < 2 {
		return false
	}

	return true
}

func cleanupTag(response string) (string, bool) {
	// Basic cleanup
	tag := strings.TrimSpace(response)
	tag = strings.Trim(tag, "\"'`.,;:") // Remove common punctuation

	// Remove numbering prefixes like "1. " or "1) "
	if idx := strings.IndexAny(tag, ".)"); idx > 0 {
		if num := strings.TrimSpace(tag[:idx]); len(num) <= 2 {
			if _, err := strconv.Atoi(num); err == nil {
				tag = strings.TrimSpace(tag[idx+1:])
			}
		}
	}

	parts := strings.Fields(tag)
	if len(parts) > 0 {
		return "", false
	}

	tag = parts[0]

	if isValidTag(tag) {
		return tag, true
	}
	return "", false
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

		if tag, valid := cleanupTag(response); valid {
			fmt.Printf("Got response from ollama (approved): %s\n", response)
			return tag, nil
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
