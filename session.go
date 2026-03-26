package biometry

import (
	"context"
	"net/http"
)

// SessionService handles biometric session lifecycle.
type SessionService struct{ c *Client }

// Start begins a new biometric session and returns the session ID.
func (s *SessionService) Start(ctx context.Context, opts ...CallOption) (string, error) {
	var id string
	req, err := s.c.newRequest(ctx, http.MethodPost, pathGateway+"/sessions/start", nil)
	if err != nil {
		return "", err
	}
	applyCallOptions(req, opts)
	if err := s.c.do(req, &id); err != nil {
		return "", err
	}
	return id, nil
}

// End terminates a session by ID.
func (s *SessionService) End(ctx context.Context, id string, req *EndSessionRequest, opts ...CallOption) error {
	return s.c.doJSON(ctx, http.MethodPost, pathGateway+"/sessions/end/"+id, req, nil, opts...)
}
