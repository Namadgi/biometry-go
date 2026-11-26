package biometry

import (
	"context"
	"net/http"
	"net/url"
)

// SessionService handles biometric session lifecycle.
type SessionService struct{ c *Client }

// Start begins a new biometric session and returns the session ID.
func (s *SessionService) Start(ctx context.Context, opts ...CallOption) (string, error) {
	req, err := s.c.newRequest(ctx, http.MethodPost, pathGatewayV2+"/sessions", nil)
	if err != nil {
		return "", err
	}
	applyCallOptions(req, opts)
	var out struct {
		SessionID string `json:"session_id"`
	}
	if err := s.c.do(req, &out); err != nil {
		return "", err
	}
	return out.SessionID, nil
}

// End terminates a session by ID. Pass a non-nil req to run a SIM-swap check.
func (s *SessionService) End(ctx context.Context, id string, req *EndSessionRequest, opts ...CallOption) error {
	return s.c.doJSON(ctx, http.MethodPost, pathGatewayV2+"/sessions/"+url.PathEscape(id)+"/end", req, nil, opts...)
}
