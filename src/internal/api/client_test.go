package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreateInvoice(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/invoices" {
			t.Errorf("expected /api/v1/invoices, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or incorrect authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json")
		}

		// Decode request body
		var req InvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Amount != 10000 {
			t.Errorf("expected amount 10000, got %d", req.Amount)
		}

		// Return mock response
		resp := InvoiceResponse{
			ID:          "inv_123",
			PaymentLink: "https://pay.stripe.com/inv_123",
			PDFUrl:      "https://example.com/invoice.pdf",
			EmailSent:   false,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with mock server
	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Test CreateInvoice
	resp, err := client.CreateInvoice(InvoiceRequest{
		Amount:      10000,
		Currency:    "usd",
		ClientName:  "Test Client",
		Description: "Test invoice",
	})

	if err != nil {
		t.Fatalf("CreateInvoice() error = %v", err)
	}

	if resp.ID != "inv_123" {
		t.Errorf("expected ID inv_123, got %s", resp.ID)
	}
	if resp.PaymentLink != "https://pay.stripe.com/inv_123" {
		t.Errorf("unexpected payment link: %s", resp.PaymentLink)
	}
}

func TestCreateInvoiceError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "invalid_request",
			Message: "Amount must be positive",
		})
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	_, err := client.CreateInvoice(InvoiceRequest{
		Amount: -100,
	})

	if err == nil {
		t.Fatal("expected error for invalid request")
	}
	if err.Error() != "Amount must be positive" {
		t.Errorf("expected 'Amount must be positive', got %q", err.Error())
	}
}

func TestGetCurrentUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/me" {
			t.Errorf("expected /api/v1/me, got %s", r.URL.Path)
		}

		resp := UserResponse{
			ID:    "user_123",
			Name:  "Test User",
			Email: "test@example.com",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		t.Fatalf("GetCurrentUser() error = %v", err)
	}

	if user.ID != "user_123" {
		t.Errorf("expected ID user_123, got %s", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}
}

func TestNewClientRequiresAuth(t *testing.T) {
	// Ensure no token is available
	originalToken := os.Getenv("LANE_TOKEN")
	os.Unsetenv("LANE_TOKEN")
	defer os.Setenv("LANE_TOKEN", originalToken)

	// Use a temp home with no token file
	originalHome := os.Getenv("HOME")
	tmpDir, _ := os.MkdirTemp("", "lane-test-*")
	os.Setenv("HOME", tmpDir)
	defer func() {
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tmpDir)
	}()

	_, err := NewClient()
	if err == nil {
		t.Error("NewClient() should fail when not authenticated")
	}
}
