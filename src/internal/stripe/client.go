package stripe

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentlink"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
)

const (
	EnvStripeKey = "STRIPE_KEY"
)

// Client wraps Stripe API operations
type Client struct {
	initialized bool
}

// InvoiceRequest contains the parameters for creating an invoice
type InvoiceRequest struct {
	AmountCents int64
	ClientName  string
	Description string
	Currency    string
}

// InvoiceResult contains the result of invoice creation
type InvoiceResult struct {
	PaymentLink string
	ProductID   string
	PriceID     string
}

// NewClient creates a new Stripe client
func NewClient() (*Client, error) {
	key := os.Getenv(EnvStripeKey)
	if key == "" {
		return nil, fmt.Errorf("environment variable %s is not set", EnvStripeKey)
	}

	stripe.Key = key

	return &Client{
		initialized: true,
	}, nil
}

// CreateInvoice creates a product, price, and payment link in Stripe
func (c *Client) CreateInvoice(req InvoiceRequest) (*InvoiceResult, error) {
	if !c.initialized {
		return nil, fmt.Errorf("stripe client not initialized")
	}

	// Default currency to USD if not specified
	currency := req.Currency
	if currency == "" {
		currency = "usd"
	}

	// Build product name
	productName := req.Description
	if req.ClientName != "" {
		productName = fmt.Sprintf("%s - %s", req.ClientName, req.Description)
	}

	// Create the product
	productParams := &stripe.ProductParams{
		Name: stripe.String(productName),
	}
	if req.Description != "" {
		productParams.Description = stripe.String(req.Description)
	}

	prod, err := product.New(productParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Create the price
	priceParams := &stripe.PriceParams{
		Product:    stripe.String(prod.ID),
		UnitAmount: stripe.Int64(req.AmountCents),
		Currency:   stripe.String(currency),
	}

	pr, err := price.New(priceParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create price: %w", err)
	}

	// Create the payment link
	linkParams := &stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Price:    stripe.String(pr.ID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	link, err := paymentlink.New(linkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	return &InvoiceResult{
		PaymentLink: link.URL,
		ProductID:   prod.ID,
		PriceID:     pr.ID,
	}, nil
}
