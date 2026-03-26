package biometry

import (
	"encoding/json"
	"net/http"
)

// Response is the full API response envelope. Used for endpoints like
// process-video and match-faces that return scoring data alongside the result.
type Response struct {
	Data            json.RawMessage `json:"data,omitempty"`
	ScoringResult   json.RawMessage `json:"scoring_result,omitempty"`
	Score           float64         `json:"score,omitempty"`
	DecisionReasons []string        `json:"decision_reasons,omitempty"`
	Message         string          `json:"message,omitempty"`

	// Header contains the raw HTTP response headers for accessing
	// X-Request-ID, X-Auto-Enroll, X-Removed-Services, etc.
	Header http.Header `json:"-"`
}

// Decode unmarshals the Data field into the provided target.
func (r *Response) Decode(v any) error {
	return json.Unmarshal(r.Data, v)
}

// RequestID returns the X-Request-ID response header value.
func (r *Response) RequestID() string {
	return r.Header.Get("X-Request-Id")
}

// AutoEnrolled returns true if the X-Auto-Enroll response header is "true",
// indicating that automatic face/voice enrollment was performed.
func (r *Response) AutoEnrolled() bool {
	return r.Header.Get("X-Auto-Enroll") == "true"
}

// RemovedServices returns true if the X-Removed-Services response header is
// "true", indicating some services were dropped due to missing consent.
func (r *Response) RemovedServices() bool {
	return r.Header.Get("X-Removed-Services") == "true"
}
