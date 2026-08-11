// Package biometry provides a Go SDK for the Biometry platform APIs.
//
// The SDK targets the Biometry API v2 gateway (/api-gateway/v2). Initialize a
// client with your API key and call methods on the service sub-clients.
//
//	client := biometry.NewClient("your-api-key-jwt")
//
//	sessionID, err := client.Sessions.Start(ctx)
//
// Biometric operations identify the end user with an opaque, customer-provided
// user_id supplied on each call's input struct (not a header).
package biometry

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default Biometry API base URL.
	DefaultBaseURL = "https://api.biometrysolutions.com"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	userAgentValue = "biometry-go/0.2.0"
)

// Service path prefixes.
const (
	pathGateway      = "/api-gateway"
	pathGatewayV2    = "/api-gateway/v2"
	pathTransactions = "/api-transactions"
	pathProjects     = "/api-projects"
)

// Client is the root Biometry SDK client. It exposes service-specific
// sub-clients as public fields.
type Client struct {
	Sessions     *SessionService
	Verification *VerificationService
	Enrollment   *EnrollmentService
	FaceMatch    *FaceMatchService
	DocAuth      *DocAuthService
	Consents     *ConsentService
	Images       *ImageService
	Deepfake     *DeepfakeService
	Transactions *TransactionService
	APIKeys      *APIKeyService

	httpClient *http.Client
	baseURL    string
	apiKey     string
	userAgent  string
}

// NewClient creates a new Biometry SDK client authenticated with the given
// API key JWT. Use functional options to customize the client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: DefaultTimeout},
		baseURL:    DefaultBaseURL,
		apiKey:     apiKey,
		userAgent:  userAgentValue,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Sessions = &SessionService{c: c}
	c.Verification = &VerificationService{c: c}
	c.Enrollment = &EnrollmentService{c: c}
	c.FaceMatch = &FaceMatchService{c: c}
	c.DocAuth = &DocAuthService{c: c}
	c.Consents = &ConsentService{c: c}
	c.Images = &ImageService{c: c}
	c.Deepfake = &DeepfakeService{c: c}
	c.Transactions = &TransactionService{c: c}
	c.APIKeys = &APIKeyService{c: c}
	return c
}
