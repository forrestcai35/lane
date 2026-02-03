package ui

import (
	"strings"
	"testing"
)

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		cents    int64
		contains string
	}{
		{10000, "$100"},
		{9999, "$99.99"},
		{50, "$0.50"},
		{100, "$1"},
		{12345, "$123.45"},
	}

	for _, tt := range tests {
		t.Run(tt.contains, func(t *testing.T) {
			got := FormatAmount(tt.cents)
			// The output includes ANSI codes, so just check the value is present
			if !strings.Contains(got, tt.contains) {
				t.Errorf("FormatAmount(%d) = %q, want to contain %q", tt.cents, got, tt.contains)
			}
		})
	}
}

func TestFormatLabel(t *testing.T) {
	got := FormatLabel("Client", "Acme Corp")
	if !strings.Contains(got, "Client") {
		t.Errorf("FormatLabel() missing label")
	}
	if !strings.Contains(got, "Acme Corp") {
		t.Errorf("FormatLabel() missing value")
	}
}

func TestFormatSuccess(t *testing.T) {
	got := FormatSuccess("Invoice created")
	if !strings.Contains(got, "✓") {
		t.Errorf("FormatSuccess() missing checkmark")
	}
	if !strings.Contains(got, "Invoice created") {
		t.Errorf("FormatSuccess() missing message")
	}
}

func TestFormatError(t *testing.T) {
	got := FormatError("Something went wrong")
	if !strings.Contains(got, "✗") {
		t.Errorf("FormatError() missing X mark")
	}
	if !strings.Contains(got, "Something went wrong") {
		t.Errorf("FormatError() missing message")
	}
}

func TestFormatStep(t *testing.T) {
	got := FormatStep("Connecting...")
	if !strings.Contains(got, "→") {
		t.Errorf("FormatStep() missing arrow")
	}
	if !strings.Contains(got, "Connecting...") {
		t.Errorf("FormatStep() missing message")
	}
}

func TestFormatLink(t *testing.T) {
	url := "https://example.com/invoice/123"
	got := FormatLink(url)
	if !strings.Contains(got, url) {
		t.Errorf("FormatLink() missing URL")
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{100.0, "100"},
		{99.99, "99.99"},
		{0.5, "0.50"},
		{1.1, "1.10"},
		{1000.0, "1000"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := formatFloat(tt.input)
			if got != tt.expected {
				t.Errorf("formatFloat(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
