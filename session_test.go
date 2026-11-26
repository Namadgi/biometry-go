package biometry

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestSessionStart(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/sessions") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"session_id": "session-abc"},
			"meta": map[string]string{},
		})
	})

	id, err := client.Sessions.Start(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "session-abc" {
		t.Errorf("expected session-abc, got %s", id)
	}
}

func TestSessionEnd(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/sessions/sess-123/end") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"meta":{"message":"ok"}}`))
	})

	err := client.Sessions.End(context.Background(), "sess-123", &EndSessionRequest{
		PhoneNumber: "+1234567890",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
