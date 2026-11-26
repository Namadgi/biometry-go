package biometry

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// parseRequestPart parses the JSON "request" multipart field from an incoming
// request into a generic map.
func parseRequestPart(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		t.Fatalf("parse multipart: %v", err)
	}
	raw := r.FormValue("request")
	if raw == "" {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("decode request part: %v", err)
	}
	return out
}

func okEnvelope(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(map[string]any{"data": data, "meta": map[string]string{}})
}

func TestEnrollFace_RequestPart(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/enrollments/face") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		req := parseRequestPart(t, r)
		if req["user_id"] != "user-123" {
			t.Errorf("expected user_id user-123, got %v", req["user_id"])
		}
		if req["is_document"] != true {
			t.Errorf("expected is_document true, got %v", req["is_document"])
		}
		if _, _, err := r.FormFile("face"); err != nil {
			t.Errorf("missing face file: %v", err)
		}
		okEnvelope(w, map[string]any{"enrolled": true, "user_id": "user-123"})
	})

	_, err := client.Enrollment.Face(context.Background(), EnrollFaceInput{
		UserID: "user-123", Face: strings.NewReader("img"), FileName: "f.jpg", IsDocument: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLiveness_ExcludeServices(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/liveness") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		req := parseRequestPart(t, r)
		if req["phrase"] != "one two" {
			t.Errorf("expected phrase, got %v", req["phrase"])
		}
		svc, ok := req["services"].(map[string]any)
		if !ok {
			t.Fatalf("expected services object, got %v", req["services"])
		}
		excl, ok := svc["exclude"].([]any)
		if !ok || len(excl) != 1 || excl[0] != ServiceVoiceRecognition {
			t.Errorf("unexpected exclude: %v", svc["exclude"])
		}
		okEnvelope(w, map[string]any{"face_liveness_detection": map[string]any{}})
	})

	_, err := client.Verification.Liveness(context.Background(), LivenessInput{
		UserID: "user-123", Video: strings.NewReader("vid"), FileName: "v.mp4",
		Phrase: "one two", ExcludeServices: []string{ServiceVoiceRecognition},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFaceVerify_UseSessionVideo(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/face-verify") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		req := parseRequestPart(t, r)
		if req["use_session_video"] != true {
			t.Errorf("expected use_session_video true, got %v", req["use_session_video"])
		}
		if req["session_id"] != "sess-1" {
			t.Errorf("expected session_id sess-1, got %v", req["session_id"])
		}
		if _, _, err := r.FormFile("video"); err == nil {
			t.Error("expected no video file when reusing session video")
		}
		okEnvelope(w, map[string]any{"face_recognition": map[string]any{}})
	})

	_, err := client.Verification.FaceVerify(context.Background(), VerifyInput{
		UserID: "user-123", Phrase: "one two", SessionID: "sess-1", UseSessionVideo: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFaceMatch_ReferenceAndVideo(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/face-match") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		req := parseRequestPart(t, r)
		if req["user_id"] != "user-123" {
			t.Errorf("expected user_id, got %v", req["user_id"])
		}
		if _, _, err := r.FormFile("reference_image"); err != nil {
			t.Errorf("missing reference_image: %v", err)
		}
		if _, _, err := r.FormFile("video"); err != nil {
			t.Errorf("missing video: %v", err)
		}
		okEnvelope(w, map[string]any{"face_recognition": map[string]any{}})
	})

	_, err := client.FaceMatch.Match(context.Background(), FaceMatchInput{
		UserID:         "user-123",
		ReferenceImage: strings.NewReader("ref"), ReferenceName: "ref.jpg",
		Video: strings.NewReader("vid"), VideoName: "v.mp4",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDocAuth_ProviderFields(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/documents/check") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		req := parseRequestPart(t, r)
		if req["provider"] != "idscan" || req["mrz_provider"] != "inhouse" {
			t.Errorf("unexpected provider fields: %v / %v", req["provider"], req["mrz_provider"])
		}
		okEnvelope(w, map[string]any{"first_name": "Jane"})
	})

	_, err := client.DocAuth.Check(context.Background(), DocAuthInput{
		Document: strings.NewReader("doc"), FileName: "p.jpg",
		Provider: "idscan", MRZProvider: "inhouse",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReferenceImageUpload_IsPut(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/users/user-123/reference-image") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if _, _, err := r.FormFile("image"); err != nil {
			t.Errorf("missing image file: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"meta":{}}`))
	})

	err := client.Images.Upload(context.Background(), "user-123", strings.NewReader("img"), "ref.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeepfakeSubmit_RequestPart(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api-gateway/v2/deepfake/checks") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		req := parseRequestPart(t, r)
		if req["user_id"] != "user-123" {
			t.Errorf("expected user_id, got %v", req["user_id"])
		}
		w.WriteHeader(http.StatusCreated)
		okEnvelope(w, map[string]any{"id": "chk-1", "status": "pending"})
	})

	check, err := client.Deepfake.Submit(context.Background(), DeepfakeInput{
		UserID: "user-123", Video: strings.NewReader("vid"), FileName: "v.mp4",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if check.ID != "chk-1" || check.Status != "pending" {
		t.Errorf("unexpected check: %+v", check)
	}
}
