package biometry

import (
	"context"
	"io"
)

// VideoService handles biometric video processing.
type VideoService struct{ c *Client }

// ProcessVideoInput is the input for processing a biometric video.
type ProcessVideoInput struct {
	Video      io.Reader
	VideoName  string
	Phrase     string
	Vocabulary string // optional
	Trigger    string // optional
}

// Process sends a video for biometric verification. Returns the full Response
// envelope which includes scoring results, score, and decision reasons.
func (v *VideoService) Process(ctx context.Context, input ProcessVideoInput, opts ...CallOption) (*Response, error) {
	fields := map[string]string{"phrase": input.Phrase}
	if input.Vocabulary != "" {
		fields["vocabulary"] = input.Vocabulary
	}
	if input.Trigger != "" {
		fields["trigger"] = input.Trigger
	}
	files := []FileUpload{{
		FieldName: "video",
		FileName:  input.VideoName,
		Content:   input.Video,
	}}
	return v.c.doMultipartRaw(ctx, pathGateway+"/process-video", fields, files, opts...)
}
