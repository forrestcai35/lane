package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/forrestcai35/lane/internal/config"
	"github.com/forrestcai35/lane/internal/ui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Lane",
	Long: `Log in to your Lane account.

This will open your browser to authenticate. Once complete,
your CLI will be automatically connected.`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

type cliAuthResponse struct {
	Code   string `json:"code,omitempty"`
	Token  string `json:"token,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println(ui.Logo.Render("⚡ Lane Login"))
	fmt.Println()

	// Check if already logged in
	if config.IsLoggedIn() {
		fmt.Println(ui.Subtle.Render("Already logged in. Use 'lane logout' to switch accounts."))
		return nil
	}

	apiURL := config.GetAPIURL()

	// Step 1: Create pending auth session
	fmt.Println(ui.FormatStep("Starting authentication..."))
	fmt.Println(ui.Subtle.Render("API: " + apiURL))

	resp, err := http.Post(apiURL+"/api/auth/cli", "application/json", nil)
	if err != nil {
		fmt.Println(ui.FormatError("Failed to connect to Lane API: " + err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(ui.FormatError(fmt.Sprintf("API returned status %d", resp.StatusCode)))
		return fmt.Errorf("API error: %s", resp.Status)
	}

	var authResp cliAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		fmt.Println(ui.FormatError("Failed to parse response: " + err.Error()))
		return err
	}

	if authResp.Code == "" {
		if authResp.Error != "" {
			fmt.Println(ui.FormatError("API error: " + authResp.Error))
		} else {
			fmt.Println(ui.FormatError("Failed to start authentication - no code received"))
		}
		return fmt.Errorf("no auth code received")
	}

	// Step 2: Open browser
	authURL := fmt.Sprintf("%s/auth/cli?code=%s", apiURL, authResp.Code)
	fmt.Println(ui.FormatStep("Opening browser..."))
	fmt.Println(ui.Subtle.Render(authURL))
	fmt.Println()

	if err := openBrowser(authURL); err != nil {
		fmt.Println(ui.Subtle.Render("Could not open browser. Please visit the URL above."))
	}

	// Step 3: Poll for token
	fmt.Println(ui.FormatStep("Waiting for authentication..."))
	fmt.Println(ui.Subtle.Render("Complete the login in your browser."))
	fmt.Println()

	token, err := pollForToken(apiURL, authResp.Code)
	if err != nil {
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Step 4: Save token
	if err := config.SaveAuthToken(token); err != nil {
		fmt.Println(ui.FormatError("Failed to save token"))
		return err
	}

	fmt.Println(ui.FormatSuccess("✓ Logged in successfully!"))
	fmt.Println()
	fmt.Println(ui.Subtle.Render("You can now create invoices with:"))
	fmt.Println(ui.Label.Render("  lane 100 --client \"Acme\" --desc \"Consulting\""))
	fmt.Println()

	return nil
}

func pollForToken(apiURL, code string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	pollURL := fmt.Sprintf("%s/api/auth/cli?code=%s", apiURL, code)

	// Poll for up to 5 minutes
	deadline := time.Now().Add(5 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				return "", fmt.Errorf("authentication timed out")
			}

			resp, err := client.Get(pollURL)
			if err != nil {
				continue // Retry on network error
			}

			var authResp cliAuthResponse
			if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			if authResp.Token != "" {
				return authResp.Token, nil
			}

			if authResp.Error != "" {
				return "", fmt.Errorf(authResp.Error)
			}

			// Status is "pending", keep polling
		}
	}
}

// openBrowser opens the default browser to the given URL
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
