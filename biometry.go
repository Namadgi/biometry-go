// Package biometry provides a Go SDK for the Biometry platform APIs.
//
// The SDK covers all API-key-authenticated endpoints across the gateway,
// transactions, and projects services. Initialize a client with your API key
// and call methods on the service sub-clients.
//
//	client := biometry.NewClient("your-api-key-jwt",
//	    biometry.WithUserFullname("John Doe"),
//	)
//
//	sessionID, err := client.Sessions.Start(ctx)
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

	userAgentValue = "biometry-go/0.1.0"
)

// Service path prefixes.
const (
	pathGateway      = "/api-gateway"
	pathTransactions = "/api-transactions"
	pathProjects     = "/api-projects"
)

// Client is the root Biometry SDK client. It exposes service-specific
// sub-clients as public fields.
type Client struct {
	Sessions     *SessionService
	Video        *VideoService
	Enrollment   *EnrollmentService
	FaceMatch    *FaceMatchService
	DocAuth      *DocAuthService
	Consents     *ConsentService
	Images       *ImageService
	AVF          *AVFService
	Transactions *TransactionService
	APIKeys      *APIKeyService

	httpClient   *http.Client
	baseURL      string
	apiKey       string
	userFullname string
	userAgent    string
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
	c.Video = &VideoService{c: c}
	c.Enrollment = &EnrollmentService{c: c}
	c.FaceMatch = &FaceMatchService{c: c}
	c.DocAuth = &DocAuthService{c: c}
	c.Consents = &ConsentService{c: c}
	c.Images = &ImageService{c: c}
	c.AVF = &AVFService{c: c}
	c.Transactions = &TransactionService{c: c}
	c.APIKeys = &APIKeyService{c: c}
	return c
}
