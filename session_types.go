package biometry

// EndSessionRequest is the request body for ending a session.
type EndSessionRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"`
}
