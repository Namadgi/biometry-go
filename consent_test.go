package biometry

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

// consentServer records the method, escaped path and JSON body of the single
// request it receives and answers with the given consent state.
func consentServer(t *testing.T, state UserConsents) (*Client, *http.Request, *map[string]any) {
	t.Helper()
	var got http.Request
	body := map[string]any{}
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		got = *r
		if r.Body != nil && r.Method == http.MethodPut {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode body: %v", err)
			}
		}
		okEnvelope(w, state)
	})
	return client, &got, &body
}

func TestConsentGive_Both(t *testing.T) {
	client, req, body := consentServer(t, UserConsents{UserID: "user-123", Authorization: true, Storage: true})

	got, err := client.Consents.Give(context.Background(), "user-123", ConsentChange{Authorization: true, Storage: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodPut || req.URL.Path != "/api-gateway/v2/users/user-123/consents" {
		t.Errorf("request = %s %s", req.Method, req.URL.Path)
	}
	want := map[string]any{"authorization": true, "storage": true}
	if len(*body) != len(want) || (*body)["authorization"] != true || (*body)["storage"] != true {
		t.Errorf("body = %v, want %v", *body, want)
	}
	if !got.Authorization || !got.Storage || got.UserID != "user-123" {
		t.Errorf("state = %+v", got)
	}
}

func TestConsentRevoke_StorageOnlyWithSession(t *testing.T) {
	client, req, body := consentServer(t, UserConsents{UserID: "user-123", Authorization: true, Storage: false})

	got, err := client.Consents.Revoke(context.Background(), "user-123", ConsentChange{Storage: true, SessionID: "sess_1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodPut {
		t.Errorf("method = %s, want PUT", req.Method)
	}
	if _, ok := (*body)["authorization"]; ok {
		t.Errorf("authorization must be omitted (left unchanged), body = %v", *body)
	}
	if (*body)["storage"] != false || (*body)["session_id"] != "sess_1" {
		t.Errorf("body = %v", *body)
	}
	if !got.Authorization || got.Storage {
		t.Errorf("state = %+v", got)
	}
}

func TestConsentGive_NothingSelected(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request expected")
	})
	for name, call := range map[string]func() (*UserConsents, error){
		"give": func() (*UserConsents, error) {
			return client.Consents.Give(context.Background(), "user-123", ConsentChange{})
		},
		"revoke": func() (*UserConsents, error) {
			return client.Consents.Revoke(context.Background(), "user-123", ConsentChange{SessionID: "s"})
		},
	} {
		if _, err := call(); !errors.Is(err, ErrNoConsentSelected) {
			t.Errorf("%s: err = %v, want ErrNoConsentSelected", name, err)
		}
	}
}

func TestConsentStatus(t *testing.T) {
	client, req, _ := consentServer(t, UserConsents{UserID: "a/b", Authorization: true})

	got, err := client.Consents.Status(context.Background(), "a/b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodGet || req.URL.EscapedPath() != "/api-gateway/v2/users/a%2Fb/consents" {
		t.Errorf("request = %s %s", req.Method, req.URL.EscapedPath())
	}
	if !got.Authorization || got.Storage {
		t.Errorf("state = %+v", got)
	}
}

func TestConsentGive_APIError(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":{"code":"session_closed","message":"session is closed"}}`))
	})
	_, err := client.Consents.Give(context.Background(), "user-123", ConsentChange{Authorization: true, Storage: true, SessionID: "sess_1"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict || apiErr.Code != "session_closed" {
		t.Errorf("err = %v, want 409 session_closed APIError", err)
	}
}
