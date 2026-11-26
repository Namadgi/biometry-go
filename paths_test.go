package biometry

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPathSegmentsAreEscaped(t *testing.T) {
	var got string
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.EscapedPath()
		okEnvelope(w, map[string]any{"id": "r", "status": "pending"})
	})
	ctx := context.Background()

	cases := []struct {
		name string
		call func() error
		want string
	}{
		{"Sessions.End", func() error { return client.Sessions.End(ctx, "a/b?c", &EndSessionRequest{}) }, "/api-gateway/v2/sessions/a%2Fb%3Fc/end"},
		{"Deepfake.Status", func() error { _, err := client.Deepfake.Status(ctx, "r#1"); return err }, "/api-gateway/v2/deepfake/checks/r%231"},
		{"Consents.Approve", func() error { return client.Consents.Approve(ctx, "t/1", "u?1") }, "/api-gateway/v2/consents/t%2F1/approve/u%3F1"},
		{"Images.Get", func() error { _, err := client.Images.Get(ctx, "u/1"); return err }, "/api-gateway/v2/users/u%2F1/reference-image"},
	}
	for _, tc := range cases {
		got = ""
		_ = tc.call() // only the request path matters here
		if got != tc.want {
			t.Errorf("%s: path = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestUseSessionVideoDoesNotUploadMedia(t *testing.T) {
	var body string
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		okEnvelope(w, map[string]any{})
	})
	ctx := context.Background()

	if _, err := client.Verification.FaceVerify(ctx, VerifyInput{
		UserID: "u-1", SessionID: "sess_1", UseSessionVideo: true,
		Media: strings.NewReader("vid"), FileName: "v.mp4",
	}); err != nil {
		t.Fatalf("FaceVerify: %v", err)
	}
	if strings.Contains(body, `name="video"`) {
		t.Error("FaceVerify uploaded media despite UseSessionVideo")
	}

	if _, err := client.FaceMatch.Match(ctx, FaceMatchInput{
		UserID: "u-1", ReferenceImage: strings.NewReader("img"), ReferenceName: "r.jpg",
		SessionID: "sess_1", UseSessionVideo: true, Video: strings.NewReader("vid"), VideoName: "v.mp4",
	}); err != nil {
		t.Fatalf("Match: %v", err)
	}
	if strings.Contains(body, `name="video"`) || !strings.Contains(body, `name="reference_image"`) {
		t.Errorf("Match body has wrong parts: %q", body)
	}

	if _, err := client.Verification.FaceVerify(ctx, VerifyInput{
		UserID: "u-1", Media: strings.NewReader("vid"), FileName: "v.mp4",
	}); err != nil {
		t.Fatalf("FaceVerify: %v", err)
	}
	if !strings.Contains(body, `name="video"`) {
		t.Error("FaceVerify without UseSessionVideo must upload media")
	}
}
