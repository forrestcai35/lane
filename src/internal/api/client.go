package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/forrestcai35/lane/internal/config"
)

// Client handles communication with the Lane API
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Lane API client
func NewClient() (*Client, error) {
	token, err := config.GetAuthToken()
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: config.GetAPIURL(),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// InvoiceRequest is the request body for creating an invoice
type InvoiceRequest struct {
	Amount      int64  `json:"amount"`       // Amount in cents
	Currency    string `json:"currency"`     // e.g., "usd"
	ClientName  string `json:"client_name"`  // Client's name
	ClientEmail string `json:"client_email"` // Client's email (for sending)
	Description string `json:"description"`  // Invoice description
	SendEmail   bool   `json:"send_email"`   // Whether to send email
}

// InvoiceResponse is the response from creating an invoice
type InvoiceResponse struct {
	ID          string `json:"id"`           // Invoice ID
	PaymentLink string `json:"payment_link"` // Stripe payment link
	PDFUrl      string `json:"pdf_url"`      // URL to download PDF
	EmailSent   bool   `json:"email_sent"`   // Whether email was sent
}

// UserResponse is the response from the /me endpoint
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse is returned on API errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// CreateInvoice creates a new invoice via the Lane API
func (c *Client) CreateInvoice(req InvoiceRequest) (*InvoiceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.request("POST", "/api/v1/invoices", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, c.parseError(resp)
	}

	var result InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetCurrentUser returns the authenticated user's info
func (c *Client) GetCurrentUser() (*UserResponse, error) {
	resp, err := c.request("GET", "/api/v1/me", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

// request makes an authenticated HTTP request to the Lane API
func (c *Client) request(method, path string, body []byte) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Lane-CLI/0.1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseError extracts an error message from an API response
func (c *Client) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
		return fmt.Errorf("%s", errResp.Message)
	}

	return fmt.Errorf("API error: %s", resp.Status)
}
