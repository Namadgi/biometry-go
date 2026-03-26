package biometry

import (
	"context"
	"io"
)

// DocAuthService handles document authentication.
type DocAuthService struct{ c *Client }

// Check verifies the authenticity of a document image (jpg, jpeg, or png).
func (d *DocAuthService) Check(ctx context.Context, document io.Reader, fileName string, opts ...CallOption) (*Response, error) {
	files := []FileUpload{{
		FieldName: "document",
		FileName:  fileName,
		Content:   document,
	}}
	return d.c.doMultipartRaw(ctx, pathGateway+"/docauth/check", nil, files, opts...)
}
