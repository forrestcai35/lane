package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette - Purple/Ghostly theme
	Purple      = lipgloss.Color("#A855F7")
	LightPurple = lipgloss.Color("#C084FC")
	DarkPurple  = lipgloss.Color("#7C3AED")
	Ghost       = lipgloss.Color("#E9D5FF")
	Dim         = lipgloss.Color("#6B7280")
	Success     = lipgloss.Color("#10B981")
	Error       = lipgloss.Color("#EF4444")
	White       = lipgloss.Color("#FAFAFA")

	// Base styles
	Bold = lipgloss.NewStyle().Bold(true)

	// Brand styling
	Logo = lipgloss.NewStyle().
		Bold(true).
		Foreground(Purple).
		MarginBottom(1)

	// Status messages
	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Success)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Error)

	// Labels and values
	Label = lipgloss.NewStyle().
		Foreground(Dim)

	Value = lipgloss.NewStyle().
		Foreground(White).
		Bold(true)

	// Highlighted values (links, amounts)
	Highlight = lipgloss.NewStyle().
			Foreground(LightPurple).
			Bold(true)

	// Subtle text
	Subtle = lipgloss.NewStyle().
		Foreground(Dim).
		Italic(true)

	// Box for final output
	ResultBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DarkPurple).
			Padding(1, 2).
			MarginTop(1)

	// Spinner/progress indicator style
	Progress = lipgloss.NewStyle().
			Foreground(Purple)
)

// FormatAmount formats a dollar amount with styling
func FormatAmount(cents int64) string {
	dollars := float64(cents) / 100
	return Highlight.Render("$" + formatFloat(dollars))
}

// FormatLink formats a URL with styling
func FormatLink(url string) string {
	return Highlight.Render(url)
}

// FormatLabel formats a label with its value
func FormatLabel(label, value string) string {
	return Label.Render(label+": ") + Value.Render(value)
}

// FormatSuccess formats a success message
func FormatSuccess(msg string) string {
	return SuccessStyle.Render("✓ " + msg)
}

// FormatError formats an error message
func FormatError(msg string) string {
	return ErrorStyle.Render("✗ " + msg)
}

// FormatStep formats a progress step
func FormatStep(msg string) string {
	return Progress.Render("→ ") + msg
}

func formatFloat(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%.2f", f)
}
