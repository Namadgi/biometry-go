package biometry

import (
	"encoding/json"
	"fmt"
)

// APIError is returned for any HTTP 4xx/5xx response from the API.
type APIError struct {
	StatusCode int
	Message    string
	RawBody    []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("biometry: HTTP %d: %s", e.StatusCode, e.Message)
}

func parseAPIError(statusCode int, body []byte) *APIError {
	ae := &APIError{StatusCode: statusCode, RawBody: body}
	var errResp struct {
		Error   any    `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &errResp) == nil {
		if s, ok := errResp.Error.(string); ok {
			ae.Message = s
		} else {
			ae.Message = errResp.Message
		}
	}
	if ae.Message == "" {
		ae.Message = string(body)
	}
	return ae
}

// IsNotFound returns true if the error is a 404 API response.
func IsNotFound(err error) bool {
	if ae, ok := err.(*APIError); ok {
		return ae.StatusCode == 404
	}
	return false
}

// IsForbidden returns true if the error is a 403 API response.
func IsForbidden(err error) bool {
	if ae, ok := err.(*APIError); ok {
		return ae.StatusCode == 403
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 API response.
func IsUnauthorized(err error) bool {
	if ae, ok := err.(*APIError); ok {
		return ae.StatusCode == 401
	}
	return false
}
