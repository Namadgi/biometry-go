package biometry

import (
	"context"
	"net/http"
)

// TransactionService handles transaction queries.
type TransactionService struct{ c *Client }

// List returns transactions matching the given filters.
func (t *TransactionService) List(ctx context.Context, req TransactionListRequest, opts ...CallOption) ([]Transaction, error) {
	var txns []Transaction
	if err := t.c.doJSON(ctx, http.MethodPost, pathTransactions+"/transactions/list", req, &txns, opts...); err != nil {
		return nil, err
	}
	return txns, nil
}

// Count returns the number of transactions matching the given filters.
func (t *TransactionService) Count(ctx context.Context, req TransactionListRequest, opts ...CallOption) (int64, error) {
	var result struct {
		Count int64 `json:"count"`
	}
	if err := t.c.doJSON(ctx, http.MethodPost, pathTransactions+"/transactions/count", req, &result, opts...); err != nil {
		return 0, err
	}
	return result.Count, nil
}

// GetBySession returns all transactions for a session.
func (t *TransactionService) GetBySession(ctx context.Context, sessionID string, opts ...CallOption) ([]Transaction, error) {
	var txns []Transaction
	r, err := t.c.newRequest(ctx, http.MethodGet, pathTransactions+"/transactions/sessions/"+sessionID, nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(r, opts)
	if err := t.c.do(r, &txns); err != nil {
		return nil, err
	}
	return txns, nil
}

// GetSamples returns sample data for a transaction.
func (t *TransactionService) GetSamples(ctx context.Context, id string, opts ...CallOption) (any, error) {
	var result any
	r, err := t.c.newRequest(ctx, http.MethodGet, pathTransactions+"/transactions/"+id+"/samples", nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(r, opts)
	if err := t.c.do(r, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateTags updates the tags on a transaction.
func (t *TransactionService) UpdateTags(ctx context.Context, id string, tags []string, opts ...CallOption) error {
	body := struct {
		Tags []string `json:"tags"`
	}{Tags: tags}
	return t.c.doJSON(ctx, http.MethodPut, pathTransactions+"/transactions/"+id+"/tags", body, nil, opts...)
}

// Flag sets the fraud flag on a transaction.
func (t *TransactionService) Flag(ctx context.Context, id string, flag string, opts ...CallOption) error {
	body := struct {
		Flag string `json:"flag"`
	}{Flag: flag}
	return t.c.doJSON(ctx, http.MethodPut, pathTransactions+"/transactions/"+id+"/flag", body, nil, opts...)
}
