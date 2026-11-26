package biometry

import (
	"context"
	"io"
)

// FaceMatchService handles face matching against a client-provided reference
// image.
type FaceMatchService struct{ c *Client }

// FaceMatchInput is the input for a face-match request.
type FaceMatchInput struct {
	UserID         string // opaque, customer-provided identity key
	ReferenceImage io.Reader
	ReferenceName  string
	// Video is the live capture to match against the reference. May be nil when
	// UseSessionVideo is true.
	Video     io.Reader
	VideoName string
	SessionID string // optional; required when UseSessionVideo is true
	// UseSessionVideo reuses the video captured earlier in the session instead
	// of uploading Video. Requires SessionID.
	UseSessionVideo bool
}

// Match compares the live capture against the provided reference image.
func (f *FaceMatchService) Match(ctx context.Context, input FaceMatchInput, opts ...CallOption) (*Response, error) {
	fields, err := requestField(v2Request{
		UserID:          input.UserID,
		SessionID:       input.SessionID,
		UseSessionVideo: input.UseSessionVideo,
	})
	if err != nil {
		return nil, err
	}
	files := []FileUpload{{FieldName: "reference_image", FileName: input.ReferenceName, Content: input.ReferenceImage}}
	if input.Video != nil {
		files = append(files, FileUpload{FieldName: "video", FileName: input.VideoName, Content: input.Video})
	}
	return f.c.doMultipartRaw(ctx, pathGatewayV2+"/face-match", fields, files, opts...)
}
