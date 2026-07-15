package kgs_test

import (
	"testing"

	"urlshortener/internal/kgs"
)

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"zero", 0, "0"},
		{"one", 1, "1"},
		{"sixty-one", 61, "Z"},
		{"sixty-two", 62, "10"},
		{"large number", 987654321, "14Q60p"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kgs.EncodeBase62(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeBase62(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
