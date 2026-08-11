package biometry

import (
	"context"
	"io"
)

// DocAuthService handles document authentication.
type DocAuthService struct{ c *Client }

// DocAuthInput is the input for a document authenticity check.
type DocAuthInput struct {
	Document  io.Reader
	FileName  string
	SessionID string // optional
	// Provider selects the document authenticity engine: "inhouse" or "idscan".
	// Defaults to the project configuration when empty.
	Provider string
	// MRZProvider selects the MRZ extraction engine: "inhouse" or "idscan".
	// Defaults to the project configuration when empty.
	MRZProvider string
}

// Check verifies the authenticity of a document image (jpg, jpeg, or png).
func (d *DocAuthService) Check(ctx context.Context, input DocAuthInput, opts ...CallOption) (*Response, error) {
	fields, err := requestField(v2Request{
		SessionID:   input.SessionID,
		Provider:    input.Provider,
		MRZProvider: input.MRZProvider,
	})
	if err != nil {
		return nil, err
	}
	files := []FileUpload{{FieldName: "document", FileName: input.FileName, Content: input.Document}}
	return d.c.doMultipartRaw(ctx, pathGatewayV2+"/documents/check", fields, files, opts...)
}
