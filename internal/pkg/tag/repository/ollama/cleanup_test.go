package ollama

import "testing"

func TestClenup(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "1",
			input:    "   Hello  ",
			expected: "hello",
		},
		{
			name:     "2",
			input:    " Лимон   🍋   ",
			expected: "лимон",
		},
		{
			name:     "3",
			input:    " 🍋   Лимон   🍋   ",
			expected: "лимон",
		},
		{
			name:     "4",
			input:    " 🍋   LeMOn   🍋   ",
			expected: "lemon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanupTag(tt.input)
			if result != tt.expected {
				t.Errorf("For input '%s', expected '%s' but got '%s'", tt.input, tt.expected, result)
			}
		})
	}
}
