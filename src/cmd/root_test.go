package cmd

import (
	"testing"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{
			name:     "whole dollars",
			input:    "100",
			expected: 10000,
			wantErr:  false,
		},
		{
			name:     "with cents",
			input:    "99.99",
			expected: 9999,
			wantErr:  false,
		},
		{
			name:     "with dollar sign",
			input:    "$250",
			expected: 25000,
			wantErr:  false,
		},
		{
			name:     "with dollar sign and cents",
			input:    "$49.95",
			expected: 4995,
			wantErr:  false,
		},
		{
			name:     "with whitespace",
			input:    "  500  ",
			expected: 50000,
			wantErr:  false,
		},
		{
			name:     "small amount",
			input:    "0.50",
			expected: 50,
			wantErr:  false,
		},
		{
			name:     "zero amount",
			input:    "0",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "negative amount",
			input:    "-50",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "invalid string",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAmount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAmount(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("parseAmount(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
