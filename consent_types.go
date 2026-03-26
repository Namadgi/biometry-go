package biometry

import "time"

// ConsentApproval represents a user's consent approval.
type ConsentApproval struct {
	ID               string    `json:"id"`
	ConsentVersionID string    `json:"consent_version_id"`
	UserID           string    `json:"user_id"`
	IsApproved       bool      `json:"is_approved"`
	CreatedAt        time.Time `json:"created_at"`
}
