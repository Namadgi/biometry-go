package biometry

import (
	"encoding/json"
	"fmt"
)

// APIError is returned for any HTTP 4xx/5xx response from the API.
type APIError struct {
	StatusCode int
	// Code is the machine-readable error code from the API v2 error envelope
	// (e.g. "invalid_request", "face_not_detected"). Empty for legacy responses.
	Code    string
	Message string
	RawBody []byte
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("biometry: HTTP %d [%s]: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("biometry: HTTP %d: %s", e.StatusCode, e.Message)
}

func parseAPIError(statusCode int, body []byte) *APIError {
	ae := &APIError{StatusCode: statusCode, RawBody: body}

	// API v2 error envelope: {"error":{"code","message"},"meta":{...}}.
	var v2 struct {
		Error struct {
			Code    string          `json:"code"`
			Message json.RawMessage `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &v2) == nil && (v2.Error.Code != "" || len(v2.Error.Message) > 0) {
		ae.Code = v2.Error.Code
		ae.Message = decodeErrorMessage(v2.Error.Message)
	} else {
		// Legacy shapes: {"error":"..."} or {"message":"..."}.
		var legacy struct {
			Error   any    `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &legacy) == nil {
			if s, ok := legacy.Error.(string); ok {
				ae.Message = s
			} else {
				ae.Message = legacy.Message
			}
		}
	}

	if ae.Message == "" {
		ae.Message = string(body)
	}
	return ae
}

// decodeErrorMessage renders the v2 error message, which is either a string or
// a field-error array; arrays/objects are returned as their raw JSON.
func decodeErrorMessage(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
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
