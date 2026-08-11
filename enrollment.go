package biometry

import (
	"context"
	"io"
)

// EnrollmentService handles face and voice enrollment.
type EnrollmentService struct{ c *Client }

// EnrollFaceInput is the input for face enrollment.
type EnrollFaceInput struct {
	UserID     string // opaque, customer-provided identity key
	Face       io.Reader
	FileName   string
	IsDocument bool // set true if the face image is from a document
}

// Face enrolls a face image for recognition.
func (e *EnrollmentService) Face(ctx context.Context, input EnrollFaceInput, opts ...CallOption) (*Response, error) {
	fields, err := requestField(v2Request{UserID: input.UserID, IsDocument: input.IsDocument})
	if err != nil {
		return nil, err
	}
	files := []FileUpload{{
		FieldName: "face",
		FileName:  input.FileName,
		Content:   input.Face,
	}}
	return e.c.doMultipartRaw(ctx, pathGatewayV2+"/enrollments/face", fields, files, opts...)
}

// EnrollVoiceInput is the input for voice enrollment.
type EnrollVoiceInput struct {
	UserID     string // opaque, customer-provided identity key
	Voice      io.Reader
	FileName   string
	Phrase     string // the phrase spoken in the recording
	Vocabulary string // optional
}

// Voice enrolls a voice recording for recognition.
func (e *EnrollmentService) Voice(ctx context.Context, input EnrollVoiceInput, opts ...CallOption) (*Response, error) {
	fields, err := requestField(v2Request{UserID: input.UserID, Phrase: input.Phrase, Vocabulary: input.Vocabulary})
	if err != nil {
		return nil, err
	}
	files := []FileUpload{{
		FieldName: "voice",
		FileName:  input.FileName,
		Content:   input.Voice,
	}}
	return e.c.doMultipartRaw(ctx, pathGatewayV2+"/enrollments/voice", fields, files, opts...)
}
