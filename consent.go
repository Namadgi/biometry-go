package biometry

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// ConsentService records user consent (Give, Revoke, Status) and custom
// consent template approvals (Approve, MyApprovals).
type ConsentService struct{ c *Client }

// ErrNoConsentSelected is returned by Give and Revoke when the ConsentChange
// selects neither Authorization nor Storage. No request is sent.
var ErrNoConsentSelected = errors.New("biometry: select at least one of Authorization or Storage")

// Give records that the user consents to Biometry processing (Authorization)
// and/or storing (Storage) their biometric data: PUT /v2/users/{user_id}/consents.
//
// Face and voice enrollment require both consents and fail with 403 without
// them; Images.Upload requires authorization consent. Collect the user's
// agreement in your app first. Giving a consent the user already gave is a
// no-op. Returns the user's consent state after the change.
func (s *ConsentService) Give(ctx context.Context, userID string, change ConsentChange, opts ...CallOption) (*UserConsents, error) {
	return s.set(ctx, userID, change, true, opts)
}

// Revoke withdraws the user's authorization and/or storage consent. Enrollment
// fails with 403 afterwards until Give is called again. Revoking a consent the
// user never gave is a no-op. Returns the user's consent state after the change.
func (s *ConsentService) Revoke(ctx context.Context, userID string, change ConsentChange, opts ...CallOption) (*UserConsents, error) {
	return s.set(ctx, userID, change, false, opts)
}

// Status returns the user's authorization and storage consent:
// GET /v2/users/{user_id}/consents. A user who never gave consent has both false.
func (s *ConsentService) Status(ctx context.Context, userID string, opts ...CallOption) (*UserConsents, error) {
	var out UserConsents
	if err := s.c.doJSON(ctx, http.MethodGet, userConsentsPath(userID), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *ConsentService) set(ctx context.Context, userID string, change ConsentChange, given bool, opts []CallOption) (*UserConsents, error) {
	if !change.Authorization && !change.Storage {
		return nil, ErrNoConsentSelected
	}
	body := userConsentsRequest{SessionID: change.SessionID}
	if change.Authorization {
		body.Authorization = &given
	}
	if change.Storage {
		body.Storage = &given
	}
	var out UserConsents
	if err := s.c.doJSON(ctx, http.MethodPut, userConsentsPath(userID), body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func userConsentsPath(userID string) string {
	return pathGatewayV2 + "/users/" + url.PathEscape(userID) + "/consents"
}

// Approve records the user's approval of one of your project's custom consent
// templates. It does not give the authorization and storage consent that
// enrollment requires; use Give for that.
func (s *ConsentService) Approve(ctx context.Context, consentID, userID string, opts ...CallOption) error {
	req, err := s.c.newRequest(ctx, http.MethodPost, pathGatewayV2+"/consents/"+consentID+"/approve/"+userID, nil)
	if err != nil {
		return err
	}
	applyCallOptions(req, opts)
	return s.c.do(req, nil)
}

// MyApprovals returns the user's custom consent template approvals.
func (s *ConsentService) MyApprovals(ctx context.Context, userID string, opts ...CallOption) ([]ConsentApproval, error) {
	var approvals []ConsentApproval
	req, err := s.c.newRequest(ctx, http.MethodGet, pathGatewayV2+"/consents/my-approvals/"+userID, nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	if err := s.c.do(req, &approvals); err != nil {
		return nil, err
	}
	return approvals, nil
}
