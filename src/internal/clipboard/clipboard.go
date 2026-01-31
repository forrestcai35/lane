package clipboard

import (
	"github.com/atotto/clipboard"
)

// Copy copies the given text to the system clipboard.
// Returns an error if clipboard access fails.
func Copy(text string) error {
	return clipboard.WriteAll(text)
}

// IsSupported returns true if clipboard operations are supported
// on the current platform.
func IsSupported() bool {
	return !clipboard.Unsupported
}
