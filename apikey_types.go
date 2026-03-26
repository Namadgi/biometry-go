package biometry

import "time"

// APIKeyInfo contains information about an API key.
type APIKeyInfo struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	Services  []string  `json:"selected_services"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
