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
	client := NewClient("test-api-key", WithBaseURL(srv.URL))
	return client, srv
}

func TestDoJSON_Success(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"id": "123"},
			"meta": map[string]string{"message": "ok"},
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

func TestDoJSON_APIError_V2Envelope(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"code": "face_not_detected", "message": "no face detected"},
			"meta":  map[string]any{"request_id": "r9"},
		})
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("expected status 422, got %d", apiErr.StatusCode)
	}
	if apiErr.Code != "face_not_detected" {
		t.Errorf("expected code face_not_detected, got %q", apiErr.Code)
	}
	if apiErr.Message != "no face detected" {
		t.Errorf("expected 'no face detected', got %q", apiErr.Message)
	}
}

func TestDoJSON_APIError_ArrayMessage(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid_request","message":[{"field":"user_id","message":"required"}]},"meta":{}}`))
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != "invalid_request" {
		t.Errorf("expected code invalid_request, got %q", apiErr.Code)
	}
	if !strings.Contains(apiErr.Message, "user_id") {
		t.Errorf("expected field-error JSON in message, got %q", apiErr.Message)
	}
}

func TestDoJSON_APIError_LegacyShape(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad input"})
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Message != "bad input" {
		t.Errorf("expected 'bad input', got %q", apiErr.Message)
	}
}

func TestDoJSON_NilDst(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"meta":{"message":"ok"}}`))
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
			"meta": map[string]string{},
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

func TestResponseRequestID(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"id": "1"},
			"meta": map[string]string{"request_id": "req-abc-123"},
		})
	})

	resp, err := client.doMultipartRaw(
		context.Background(),
		"/test",
		nil,
		[]FileUpload{{FieldName: "f", FileName: "f.bin", Content: strings.NewReader("data")}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RequestID() != "req-abc-123" {
		t.Errorf("expected request ID req-abc-123, got %s", resp.RequestID())
	}
}

func TestRequestMetadataHeaders(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Device-Info") != `{"model":"iPhone"}` {
			t.Errorf("unexpected X-Device-Info: %s", r.Header.Get("X-Device-Info"))
		}
		if r.Header.Get("X-Request-Tags") != "vip,test" {
			t.Errorf("unexpected X-Request-Tags: %s", r.Header.Get("X-Request-Tags"))
		}
		if r.Header.Get("X-Client-App") != "demo" {
			t.Errorf("unexpected X-Client-App: %s", r.Header.Get("X-Client-App"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": "ok", "meta": map[string]string{}})
	})

	err := client.doJSON(context.Background(), http.MethodPost, "/test", nil, nil,
		WithDeviceInfo(`{"model":"iPhone"}`),
		WithRequestTags("vip,test"),
		WithClientApp("demo"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoMultipartRaw_ReturnsEnvelope(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data":     map[string]string{"id": "abc"},
			"decision": map[string]any{"status": "pass", "score": 92.5, "reasons": []string{"face_match"}},
			"meta":     map[string]string{"message": "done"},
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
	if resp.Decision == nil || resp.Decision.Status != "pass" {
		t.Fatalf("unexpected decision: %+v", resp.Decision)
	}
	if resp.Decision.Score == nil || *resp.Decision.Score != 92.5 {
		t.Errorf("expected score=92.5, got %v", resp.Decision.Score)
	}
	if len(resp.Decision.Reasons) != 1 || resp.Decision.Reasons[0] != "face_match" {
		t.Errorf("unexpected decision reasons: %v", resp.Decision.Reasons)
	}
	if resp.Meta.Message != "done" {
		t.Errorf("expected message=done, got %s", resp.Meta.Message)
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
