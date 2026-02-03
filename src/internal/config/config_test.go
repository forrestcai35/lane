package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAPIURL(t *testing.T) {
	// Save original env value
	original := os.Getenv(EnvAPIURL)
	defer os.Setenv(EnvAPIURL, original)

	t.Run("returns default when env not set", func(t *testing.T) {
		os.Unsetenv(EnvAPIURL)
		got := GetAPIURL()
		if got != DefaultAPIURL {
			t.Errorf("GetAPIURL() = %q, want %q", got, DefaultAPIURL)
		}
	})

	t.Run("returns env value when set", func(t *testing.T) {
		customURL := "https://custom.example.com"
		os.Setenv(EnvAPIURL, customURL)
		got := GetAPIURL()
		if got != customURL {
			t.Errorf("GetAPIURL() = %q, want %q", got, customURL)
		}
	})
}

func TestAuthToken(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "lane-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override HOME to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Clear any env token
	originalToken := os.Getenv(EnvAuthToken)
	os.Unsetenv(EnvAuthToken)
	defer os.Setenv(EnvAuthToken, originalToken)

	t.Run("returns error when not logged in", func(t *testing.T) {
		_, err := GetAuthToken()
		if err == nil {
			t.Error("GetAuthToken() should return error when not logged in")
		}
	})

	t.Run("save and retrieve token", func(t *testing.T) {
		testToken := "test-token-12345"
		err := SaveAuthToken(testToken)
		if err != nil {
			t.Fatalf("SaveAuthToken() error = %v", err)
		}

		got, err := GetAuthToken()
		if err != nil {
			t.Fatalf("GetAuthToken() error = %v", err)
		}
		if got != testToken {
			t.Errorf("GetAuthToken() = %q, want %q", got, testToken)
		}
	})

	t.Run("token file has correct permissions", func(t *testing.T) {
		tokenPath := filepath.Join(tmpDir, ".lane", "token")
		info, err := os.Stat(tokenPath)
		if err != nil {
			t.Fatalf("failed to stat token file: %v", err)
		}
		// Check that file is not world-readable (0600)
		perm := info.Mode().Perm()
		if perm != 0600 {
			t.Errorf("token file permissions = %o, want 0600", perm)
		}
	})

	t.Run("IsLoggedIn returns true after login", func(t *testing.T) {
		if !IsLoggedIn() {
			t.Error("IsLoggedIn() = false, want true after saving token")
		}
	})

	t.Run("delete token", func(t *testing.T) {
		err := DeleteAuthToken()
		if err != nil {
			t.Fatalf("DeleteAuthToken() error = %v", err)
		}

		if IsLoggedIn() {
			t.Error("IsLoggedIn() = true after deleting token")
		}
	})

	t.Run("env token takes precedence", func(t *testing.T) {
		envToken := "env-token-xyz"
		os.Setenv(EnvAuthToken, envToken)

		// Save a different file token
		SaveAuthToken("file-token-abc")

		got, err := GetAuthToken()
		if err != nil {
			t.Fatalf("GetAuthToken() error = %v", err)
		}
		if got != envToken {
			t.Errorf("GetAuthToken() = %q, want env token %q", got, envToken)
		}

		os.Unsetenv(EnvAuthToken)
	})
}

func TestGetConfigDir(t *testing.T) {
	dir := GetConfigDir()
	if dir == "" {
		t.Error("GetConfigDir() returned empty string")
	}
	if filepath.Base(dir) != ".lane" {
		t.Errorf("GetConfigDir() = %q, want path ending in .lane", dir)
	}
}
