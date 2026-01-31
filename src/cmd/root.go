package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/forrest/lane/internal/clipboard"
	"github.com/forrest/lane/internal/stripe"
	"github.com/forrest/lane/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Version information (set via ldflags during build)
	Version = "0.1.0"

	// Flags
	clientName  string
	description string
	currency    string
	noCopy      bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "lane <amount>",
	Short: "Generate Stripe invoices instantly",
	Long: ui.Logo.Render("Lane") + `
The fastest way to generate a Stripe invoice from the terminal.
Linear for Invoicing.

` + ui.Label.Render("Usage:") + `
  lane 500 --client "Apple" --desc "Web Design"

` + ui.Label.Render("Environment:") + `
  STRIPE_KEY    Your Stripe secret key (required)`,
	Example: `  lane 100 --client "Acme Corp" --desc "Consulting"
  lane 2500 --client "Startup Inc" --desc "Logo Design" --currency eur
  lane 50.99 --desc "Quick fix"`,
	Args: cobra.ExactArgs(1),
	RunE: runInvoice,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&clientName, "client", "c", "", "Client name")
	rootCmd.Flags().StringVarP(&description, "desc", "d", "", "Invoice description (required)")
	rootCmd.Flags().StringVar(&currency, "currency", "usd", "Currency code (usd, eur, gbp, etc.)")
	rootCmd.Flags().BoolVar(&noCopy, "no-copy", false, "Don't copy link to clipboard")

	rootCmd.MarkFlagRequired("desc")

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func runInvoice(cmd *cobra.Command, args []string) error {
	// Parse amount
	amountCents, err := parseAmount(args[0])
	if err != nil {
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Print header
	fmt.Println()
	fmt.Println(ui.Logo.Render("⚡ Lane"))

	// Initialize Stripe client
	fmt.Println(ui.FormatStep("Connecting to Stripe..."))
	client, err := stripe.NewClient()
	if err != nil {
		fmt.Println(ui.FormatError(err.Error()))
		return err
	}

	// Create the invoice
	fmt.Println(ui.FormatStep("Creating invoice..."))
	result, err := client.CreateInvoice(stripe.InvoiceRequest{
		AmountCents: amountCents,
		ClientName:  clientName,
		Description: description,
		Currency:    strings.ToLower(currency),
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

	output.WriteString(ui.FormatSuccess("Invoice created successfully!"))
	output.WriteString("\n\n")

	// Details
	if clientName != "" {
		output.WriteString(ui.FormatLabel("Client", clientName))
		output.WriteString("\n")
	}
	output.WriteString(ui.FormatLabel("Description", description))
	output.WriteString("\n")
	output.WriteString(ui.FormatLabel("Amount", ui.FormatAmount(amountCents)))
	output.WriteString("\n\n")

	// Payment link
	output.WriteString(ui.Label.Render("Payment Link: "))
	output.WriteString(clipboardStatus)
	output.WriteString("\n")
	output.WriteString(ui.FormatLink(result.PaymentLink))

	// Print the result box
	fmt.Println(ui.ResultBox.Render(output.String()))
	fmt.Println()

	return nil
}

// parseAmount converts a string amount to cents
// Accepts: "500" (dollars), "500.00" (dollars), "50.5" (dollars)
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

	// Convert to cents
	cents := int64(dollars * 100)

	return cents, nil
}
