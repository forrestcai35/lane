package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// API configuration
	DefaultAPIURL = "https://lane-website.netlify.app" // No trailing slash
	EnvAPIURL     = "LANE_API_URL"                     // Override for development
	EnvAuthToken  = "LANE_TOKEN"                       // Auth token (set by 'lane login')
)

// configDir returns the Lane config directory path
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, ".lane"), nil
}

// GetAPIURL returns the Lane API URL
func GetAPIURL() string {
	if url := os.Getenv(EnvAPIURL); url != "" {
		return url
	}
	return DefaultAPIURL
}

// GetAuthToken returns the stored auth token
func GetAuthToken() (string, error) {
	// First check env var (for CI/scripts)
	if token := os.Getenv(EnvAuthToken); token != "" {
		return token, nil
	}

	// Then check token file
	dir, err := configDir()
	if err != nil {
		return "", err
	}

	tokenPath := filepath.Join(dir, "token")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("not logged in. Run 'lane login' to authenticate")
		}
		return "", fmt.Errorf("could not read token: %w", err)
	}

	return string(data), nil
}

// SaveAuthToken stores the auth token
func SaveAuthToken(token string) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	tokenPath := filepath.Join(dir, "token")
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("could not save token: %w", err)
	}

	return nil
}

// DeleteAuthToken removes the stored auth token
func DeleteAuthToken() error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	tokenPath := filepath.Join(dir, "token")
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not delete token: %w", err)
	}

	return nil
}

// IsLoggedIn returns true if the user has a valid auth token
func IsLoggedIn() bool {
	_, err := GetAuthToken()
	return err == nil
}

// GetConfigDir returns the config directory path (for display)
func GetConfigDir() string {
	dir, _ := configDir()
	return dir
}
