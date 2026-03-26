package biometry

import (
	"context"
	"net/http"
)

// ConsentService handles custom consent approvals (API-key-authenticated only).
type ConsentService struct{ c *Client }

// Approve approves a custom consent for a user.
func (s *ConsentService) Approve(ctx context.Context, consentID, userID string, opts ...CallOption) error {
	req, err := s.c.newRequest(ctx, http.MethodPost, pathGateway+"/custom-consents/"+consentID+"/approve/"+userID, nil)
	if err != nil {
		return err
	}
	applyCallOptions(req, opts)
	return s.c.do(req, nil)
}

// MyApprovals returns the consent approvals for a user.
func (s *ConsentService) MyApprovals(ctx context.Context, userID string, opts ...CallOption) ([]ConsentApproval, error) {
	var approvals []ConsentApproval
	req, err := s.c.newRequest(ctx, http.MethodGet, pathGateway+"/custom-consents/my-approvals/"+userID, nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	if err := s.c.do(req, &approvals); err != nil {
		return nil, err
	}
	return approvals, nil
}
