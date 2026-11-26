package biometry

import (
	"encoding/json"
	"net/http"
)

// Response is the full API v2 success envelope. Used for verification endpoints
// (liveness, face-verify, voice-verify, face-match, documents/check) that may
// carry a scoring decision alongside the raw service data.
type Response struct {
	// Data is the endpoint-specific payload. Decode it with Decode.
	Data json.RawMessage `json:"data,omitempty"`
	// Decision is the scoring decision, present only when a scoring system
	// applies to the API key.
	Decision *Decision `json:"decision,omitempty"`
	// Meta carries the request ID and a human-readable message.
	Meta Meta `json:"meta"`

	// Header contains the raw HTTP response headers.
	Header http.Header `json:"-"`
}

// Decision is the scoring decision computed for the request.
type Decision struct {
	Status  string   `json:"status,omitempty"`
	Score   *float64 `json:"score,omitempty"`
	Reasons []string `json:"reasons,omitempty"`
}

// Meta is the response metadata attached to every API v2 response.
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Message   string `json:"message,omitempty"`
}

// Decode unmarshals the Data field into the provided target.
func (r *Response) Decode(v any) error {
	return json.Unmarshal(r.Data, v)
}

// RequestID returns the request identifier, preferring the envelope meta and
// falling back to the X-Request-Id response header.
func (r *Response) RequestID() string {
	if r.Meta.RequestID != "" {
		return r.Meta.RequestID
	}
	return r.Header.Get("X-Request-Id")
}
