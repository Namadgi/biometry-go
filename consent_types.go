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

// ConsentChange selects which consents Give or Revoke changes. Select at
// least one; enrollment needs both.
type ConsentChange struct {
	Authorization bool // consent to process the user's biometric data
	Storage       bool // consent to store the user's biometric data
	// SessionID optionally links the consent to a session, so it shows with
	// that session in the Console.
	SessionID string
}

// UserConsents is the user's consent state.
type UserConsents struct {
	UserID        string `json:"user_id"`
	Authorization bool   `json:"authorization"`
	Storage       bool   `json:"storage"`
}

type userConsentsRequest struct {
	Authorization *bool  `json:"authorization,omitempty"`
	Storage       *bool  `json:"storage,omitempty"`
	SessionID     string `json:"session_id,omitempty"`
}
