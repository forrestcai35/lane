package cmd

import (
	"fmt"

	"github.com/forrestcai35/lane/internal/config"
	"github.com/forrestcai35/lane/internal/ui"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of Lane",
	Long:  `Removes your stored authentication token.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		fmt.Println(ui.Subtle.Render("You're not logged in."))
		return nil
	}

	if err := config.DeleteAuthToken(); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	fmt.Println(ui.FormatSuccess("✓ Logged out"))
	return nil
}
