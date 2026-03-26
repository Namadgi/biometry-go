package biometry

import (
	"context"
	"net/http"
)

// APIKeyService manages API keys via the public endpoints.
type APIKeyService struct{ c *Client }

// Mint creates a new API key.
func (a *APIKeyService) Mint(ctx context.Context, opts ...CallOption) (*APIKeyInfo, error) {
	var info APIKeyInfo
	req, err := a.c.newRequest(ctx, http.MethodPost, pathProjects+"/public/apikeys", nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	if err := a.c.do(req, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Revoke deletes an API key by ID.
func (a *APIKeyService) Revoke(ctx context.Context, keyID string, opts ...CallOption) error {
	req, err := a.c.newRequest(ctx, http.MethodDelete, pathProjects+"/public/apikeys/"+keyID, nil)
	if err != nil {
		return err
	}
	applyCallOptions(req, opts)
	return a.c.do(req, nil)
}

// Info returns information about the current API key.
func (a *APIKeyService) Info(ctx context.Context, opts ...CallOption) (*APIKeyInfo, error) {
	var info APIKeyInfo
	req, err := a.c.newRequest(ctx, http.MethodGet, pathProjects+"/public/apikeys", nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	if err := a.c.do(req, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
