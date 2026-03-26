package biometry

import (
	"context"
	"io"
	"strings"
)

// EnrollmentService handles face and voice enrollment.
type EnrollmentService struct{ c *Client }

// EnrollFaceInput is the input for face enrollment.
type EnrollFaceInput struct {
	Face       io.Reader
	FileName   string
	IsDocument bool // set true if the face image is from a document
}

// Face enrolls a face image for recognition.
func (e *EnrollmentService) Face(ctx context.Context, input EnrollFaceInput, opts ...CallOption) (*Response, error) {
	fields := map[string]string{}
	if input.IsDocument {
		fields["is_document"] = "true"
	}
	files := []FileUpload{{
		FieldName: "face",
		FileName:  input.FileName,
		Content:   input.Face,
	}}
	return e.c.doMultipartRaw(ctx, pathGateway+"/enroll/face", fields, files, opts...)
}

// VoiceFile represents a single voice recording for enrollment.
type VoiceFile struct {
	FileName string
	Content  io.Reader
}

// EnrollVoiceInput is the input for voice enrollment.
type EnrollVoiceInput struct {
	Files      []VoiceFile
	Phrases    []string // one phrase per voice file
	Vocabulary string   // optional
}

// Voice enrolls voice recordings for recognition.
func (e *EnrollmentService) Voice(ctx context.Context, input EnrollVoiceInput, opts ...CallOption) (*Response, error) {
	fields := map[string]string{
		"phrase": strings.Join(input.Phrases, ","),
	}
	if input.Vocabulary != "" {
		fields["vocabulary"] = input.Vocabulary
	}
	files := make([]FileUpload, len(input.Files))
	for i, f := range input.Files {
		files[i] = FileUpload{
			FieldName: "voice",
			FileName:  f.FileName,
			Content:   f.Content,
		}
	}
	return e.c.doMultipartRaw(ctx, pathGateway+"/enroll/voice", fields, files, opts...)
}
