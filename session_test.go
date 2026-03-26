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
		if !strings.HasSuffix(r.URL.Path, "/sessions/start") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": "session-abc"})
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
		if !strings.HasSuffix(r.URL.Path, "/sessions/end/sess-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"ok"}`))
	})

	err := client.Sessions.End(context.Background(), "sess-123", &EndSessionRequest{
		PhoneNumber: "+1234567890",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
