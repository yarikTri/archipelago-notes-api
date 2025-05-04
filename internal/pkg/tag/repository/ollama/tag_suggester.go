package ollama

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/enescakir/emoji"

	"github.com/yarikTri/archipelago-notes-api/internal/clients/llm"
)

const (
	TagSytemPrompt = "Придумай тег, который наиболее точно описывает текст этой заметки - одно слово. Отвечай на том же языке, на котором сделан запрос."
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

// func isInvalidLLMAnswerEng(s string) bool {
// 	invalidPrefixes := []string{
// 		"i'm sorry",
// 		"i am sorry",
// 		"i apologize",
// 		"as an ai",
// 		"as a language model",
// 		"i can't",
// 		"i cannot",
// 		"unfortunately",
// 		"my purpose is",
// 		"i don't",
// 		"i do not",
// 		"this content",
// 		"this request",
// 		"i am not",
// 		"i am unable",
// 		"that's not",
// 		"that is outside",
// 	}

// 	for _, prefix := range invalidPrefixes {
// 		if strings.HasPrefix(s, prefix) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func isInvalidLLMAnswerRu(s string) bool {
// 	invalidPrefixes := []string{
// 		"я извиняюсь",
// 		"я сожалею",
// 		"как ии",
// 		"как искусственный интеллект",
// 		"я не могу",
// 		"к сожалению",
// 		"моя функция",
// 		"это выходит",
// 		"данный запрос",
// 	}

// 	for _, prefix := range invalidPrefixes {
// 		if strings.HasPrefix(s, prefix) {
// 			return true
// 		}
// 	}
// 	return false
// }

func removeAllEmojis(s string) string {
	isEmoji := func(r rune) bool {
		_, exists := emoji.Map()[string(r)]
		return exists
	}

	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		if !isEmoji(r) {
			_, err := builder.WriteRune(r)
			if err != nil {
				panic("got invalid rune")
			}
		}
	}

	return builder.String()
}

func isValidTag(tag string) bool {
	if len(tag) == 0 {
		return false
	}

	// if isInvalidLLMAnswerEng(tag) {
	// 	return false
	// }

	// if isInvalidLLMAnswerRu(tag) {
	// 	return false
	// }

	parts := strings.Fields(tag)
	if len(parts) != 1 {
		return false
	}
	tag = parts[0]

	if len(tag) > 80 || len(tag) < 2 {
		return false
	}

	return true
}

func cleanupTag(response string) string {
	tag := strings.TrimSpace(response)

	tag = removeAllEmojis(tag)

	// Compile the regular expression
	// This matches any character that is NOT a letter (a-z, A-Z) or number (0-9)
	reg := regexp.MustCompile(`[^a-zA-Z0-9]`)

	// Replace all matched characters with an empty string
	tag = reg.ReplaceAllString(tag, "")

	return strings.TrimSpace(strings.ToLower(tag))
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

		fmt.Printf("Tag response before clenup: %s\n", response)
		tag := cleanupTag(response)
		fmt.Printf("Tag after cleanup: %s\n", tag)

		if isValidTag(tag) {
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
