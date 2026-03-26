package biometry

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := NewClient("test-api-key",
		WithBaseURL(srv.URL),
		WithUserFullname("Test User"),
	)
	return client, srv
}

func TestDoJSON_Success(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-User-Fullname") != "Test User" {
			t.Errorf("unexpected fullname header: %s", r.Header.Get("X-User-Fullname"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data":    map[string]string{"id": "123"},
			"message": "ok",
		})
	})

	var result struct {
		ID string `json:"id"`
	}
	err := client.doJSON(context.Background(), http.MethodPost, "/test", map[string]string{"key": "val"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "123" {
		t.Errorf("expected id=123, got %s", result.ID)
	}
}

func TestDoJSON_APIError(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "bad input"})
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "bad input" {
		t.Errorf("expected 'bad input', got %q", apiErr.Message)
	}
}

func TestDoJSON_NilDst(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"ok"}`))
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoMultipart_Success(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("unexpected content type: %s", ct)
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if r.FormValue("field1") != "value1" {
			t.Errorf("unexpected field1: %s", r.FormValue("field1"))
		}
		file, _, err := r.FormFile("upload")
		if err != nil {
			t.Fatalf("get form file: %v", err)
		}
		defer file.Close()
		content, _ := io.ReadAll(file)
		if string(content) != "file-content" {
			t.Errorf("unexpected file content: %s", string(content))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"result": "ok"},
		})
	})

	var result struct {
		Result string `json:"result"`
	}
	err := client.doMultipart(
		context.Background(),
		"/upload",
		map[string]string{"field1": "value1"},
		[]FileUpload{{FieldName: "upload", FileName: "test.bin", Content: strings.NewReader("file-content")}},
		&result,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result != "ok" {
		t.Errorf("expected result=ok, got %s", result.Result)
	}
}

func TestDoMultipartRaw_ReturnsEnvelope(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data":             map[string]string{"id": "abc"},
			"scoring_result":   map[string]string{"status": "approved"},
			"score":            92.5,
			"decision_reasons": []string{"face_match"},
			"message":          "done",
		})
	})

	resp, err := client.doMultipartRaw(
		context.Background(),
		"/process",
		nil,
		[]FileUpload{{FieldName: "video", FileName: "test.mp4", Content: strings.NewReader("video-data")}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Score != 92.5 {
		t.Errorf("expected score=92.5, got %f", resp.Score)
	}
	if len(resp.DecisionReasons) != 1 || resp.DecisionReasons[0] != "face_match" {
		t.Errorf("unexpected decision reasons: %v", resp.DecisionReasons)
	}
	if resp.Message != "done" {
		t.Errorf("expected message=done, got %s", resp.Message)
	}
}

func TestCallOptions(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-User-Fullname") != "Override Name" {
			t.Errorf("expected overridden fullname, got %s", r.Header.Get("X-User-Fullname"))
		}
		if r.Header.Get("X-Session-ID") != "sess-123" {
			t.Errorf("expected session ID, got %s", r.Header.Get("X-Session-ID"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"id-1"}`))
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil,
		WithFullname("Override Name"),
		WithSessionID("sess-123"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestErrorPredicates(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		checkFn    func(error) bool
		expectTrue bool
	}{
		{"not found", 404, IsNotFound, true},
		{"forbidden", 403, IsForbidden, true},
		{"unauthorized", 401, IsUnauthorized, true},
		{"not found for 400", 400, IsNotFound, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.status, Message: "test"}
			if got := tt.checkFn(err); got != tt.expectTrue {
				t.Errorf("expected %v, got %v", tt.expectTrue, got)
			}
		})
	}
}
