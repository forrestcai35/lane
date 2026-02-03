package cmd

import (
	"fmt"
	"strings"

	"github.com/forrestcai35/lane/internal/api"
	"github.com/forrestcai35/lane/internal/clipboard"
	"github.com/forrestcai35/lane/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Version information (set via ldflags during build)
	Version = "0.1.0"

	// Flags
	clientName  string
	clientEmail string
	description string
	currency    string
	sendEmail   bool
	noCopy      bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "lane <amount>",
	Short: "Generate Stripe invoices instantly",
	Long: ui.Logo.Render("Lane") + `
The fastest way to generate a Stripe invoice from the terminal.

` + ui.Label.Render("Quick Start:") + `
  lane login                              # Authenticate with Lane
  lane 500 --client "Apple" --desc "Work" # Create invoice`,
	Example: `  lane 100 --client "Acme Corp" --desc "Consulting"
  lane 500 --client "Apple" --desc "Web Design" --email "tim@apple.com" --send
  lane 2500 --desc "Logo Design" --currency eur`,
	Args: cobra.ExactArgs(1),
	RunE: runInvoice,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&clientName, "client", "c", "", "Client name")
	rootCmd.Flags().StringVarP(&clientEmail, "email", "e", "", "Client email address")
	rootCmd.Flags().StringVarP(&description, "desc", "d", "", "Invoice description (required)")
	rootCmd.Flags().StringVar(&currency, "currency", "usd", "Currency code (usd, eur, gbp, etc.)")
	rootCmd.Flags().BoolVar(&sendEmail, "send", false, "Send invoice via email (requires --email)")
	rootCmd.Flags().BoolVar(&noCopy, "no-copy", false, "Don't copy link to clipboard")

	rootCmd.MarkFlagRequired("desc")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func runInvoice(cmd *cobra.Command, args []string) error {
	// Parse amount
	amountCents, err := parseAmount(args[0])
	if err != nil {
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Validate email flags
	if sendEmail && clientEmail == "" {
		err := fmt.Errorf("--send requires --email flag")
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Print header
	fmt.Println()
	fmt.Println(ui.Logo.Render("⚡ Lane"))

	// Initialize API client
	fmt.Println(ui.FormatStep("Connecting..."))
	client, err := api.NewClient()
	if err != nil {
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Create the invoice via API
	fmt.Println(ui.FormatStep("Creating invoice..."))
	result, err := client.CreateInvoice(api.InvoiceRequest{
		Amount:      amountCents,
		Currency:    strings.ToLower(currency),
		ClientName:  clientName,
		ClientEmail: clientEmail,
		Description: description,
		SendEmail:   sendEmail,
	})
	if err != nil {
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Copy to clipboard
	clipboardStatus := ""
	if !noCopy && clipboard.IsSupported() {
		if err := clipboard.Copy(result.PaymentLink); err != nil {
			clipboardStatus = ui.Subtle.Render("(clipboard unavailable)")
		} else {
			clipboardStatus = ui.SuccessStyle.Render("(copied!)")
		}
	}

	// Build output
	var output strings.Builder

	output.WriteString(ui.FormatSuccess("Invoice created!"))
	output.WriteString("\n\n")

	// Details
	if clientName != "" {
		output.WriteString(ui.FormatLabel("Client", clientName))
		output.WriteString("\n")
	}
	if clientEmail != "" {
		output.WriteString(ui.FormatLabel("Email", clientEmail))
		output.WriteString("\n")
	}
	output.WriteString(ui.FormatLabel("Description", description))
	output.WriteString("\n")
	output.WriteString(ui.FormatLabel("Amount", ui.FormatAmount(amountCents)))
	output.WriteString("\n")
	output.WriteString(ui.FormatLabel("Invoice", result.ID))
	output.WriteString("\n\n")

	// Status messages
	if result.EmailSent {
		output.WriteString(ui.FormatSuccess("✓ Email sent to " + clientEmail))
		output.WriteString("\n\n")
	}

	// Payment link
	output.WriteString(ui.Label.Render("Payment Link: "))
	output.WriteString(clipboardStatus)
	output.WriteString("\n")
	output.WriteString(ui.FormatLink(result.PaymentLink))

	fmt.Println(ui.ResultBox.Render(output.String()))
	fmt.Println()

	return nil
}

// parseAmount converts a string amount to cents
func parseAmount(s string) (int64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")

	var dollars float64
	_, err := fmt.Sscanf(s, "%f", &dollars)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %s", s)
	}

	if dollars <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}

	cents := int64(dollars * 100)
	return cents, nil
}
