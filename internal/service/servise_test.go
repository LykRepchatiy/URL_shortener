package service

import (
	"strings"
	"testing"
)

func TestShortURL(t *testing.T) {
	input := "testInput"
	result := ShortURL(input)
	if len(result) != 10 {
		t.Errorf("Тест 1: ожидалась длина 10, получено %d", len(result))
	}

	for _, char := range result {
		if !strings.ContainsRune(Alphabet, char) {
			t.Errorf("Тест 2: символ '%c' не принадлежит алфавиту", char)
		}
	}

	longInput := strings.Repeat("a", 1000)
	longResult := ShortURL(longInput)
	if len(longResult) != 10 {
		t.Errorf("Тест 3: ожидалась длина 10 для длинной строки, получено %d", len(longResult))
	}
}

func TestIsValidURL(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {"Valid HTTP URL", "http://example.com", true},
        {"Valid HTTPS URL", "https://google.com/path", true},
        {"Valid URL with Query", "https://example.com?query=123", true},
		
        {"Valid URL with Fragment", "https://example.com#fragment", false},
        {"Empty String", "", false},
        {"Missing Scheme", "example.com", false},
        {"Missing Host", "http://", false},
        {"Invalid Characters", "http://exa mple.com", false},
        {"Relative Path", "/path", false},
        {"Local Path", "./local", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsValidURL(tt.input)
            if result != tt.expected {
                t.Errorf("Тест '%s': для входа '%s' ожидалось %v, получено %v", tt.name, tt.input, tt.expected, result)
            }
        })
    }
}
