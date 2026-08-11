package biometry

import (
	"context"
	"io"
)

// VerificationService handles biometric verification. API v2 splits the v1
// process-video call into three focused operations: Liveness, FaceVerify, and
// VoiceVerify.
type VerificationService struct{ c *Client }

// LivenessInput is the input for a liveness check (face liveness, active
// speaker detection, visual speech recognition).
type LivenessInput struct {
	UserID          string // opaque, customer-provided identity key
	Video           io.Reader
	FileName        string
	Phrase          string
	Vocabulary      string   // optional
	Trigger         string   // optional contextual trigger event name
	SessionID       string   // optional, for combined scoring
	ExcludeServices []string // optional services to skip (see Service* constants)
}

// Liveness runs the liveness pipeline on the submitted video.
func (v *VerificationService) Liveness(ctx context.Context, input LivenessInput, opts ...CallOption) (*Response, error) {
	req := v2Request{
		UserID:     input.UserID,
		Phrase:     input.Phrase,
		Vocabulary: input.Vocabulary,
		Trigger:    input.Trigger,
		SessionID:  input.SessionID,
	}
	if len(input.ExcludeServices) > 0 {
		req.Services = &v2Services{Exclude: input.ExcludeServices}
	}
	fields, err := requestField(req)
	if err != nil {
		return nil, err
	}
	files := []FileUpload{{FieldName: "video", FileName: input.FileName, Content: input.Video}}
	return v.c.doMultipartRaw(ctx, pathGatewayV2+"/liveness", fields, files, opts...)
}

// VerifyInput is the input for face or voice verification against the user's
// enrolled template.
type VerifyInput struct {
	UserID string // opaque, customer-provided identity key
	// Media is the video (face) or audio/video (voice) file. May be nil when
	// UseSessionVideo is true.
	Media      io.Reader
	FileName   string
	Phrase     string
	Vocabulary string // optional
	SessionID  string // optional, for combined scoring / session-video reuse
	// UseSessionVideo reuses the video captured earlier in the session instead
	// of uploading Media. Requires SessionID.
	UseSessionVideo bool
}

// FaceVerify verifies the face from the input video/image against the user's
// enrolled face template.
func (v *VerificationService) FaceVerify(ctx context.Context, input VerifyInput, opts ...CallOption) (*Response, error) {
	return v.verify(ctx, "/face-verify", "video", input, opts)
}

// VoiceVerify verifies the voice from the input audio against the user's
// enrolled voice template.
func (v *VerificationService) VoiceVerify(ctx context.Context, input VerifyInput, opts ...CallOption) (*Response, error) {
	return v.verify(ctx, "/voice-verify", "voice", input, opts)
}

func (v *VerificationService) verify(ctx context.Context, path, fileField string, input VerifyInput, opts []CallOption) (*Response, error) {
	fields, err := requestField(v2Request{
		UserID:          input.UserID,
		Phrase:          input.Phrase,
		Vocabulary:      input.Vocabulary,
		SessionID:       input.SessionID,
		UseSessionVideo: input.UseSessionVideo,
	})
	if err != nil {
		return nil, err
	}
	var files []FileUpload
	if input.Media != nil {
		files = []FileUpload{{FieldName: fileField, FileName: input.FileName, Content: input.Media}}
	}
	return v.c.doMultipartRaw(ctx, pathGatewayV2+path, fields, files, opts...)
}
