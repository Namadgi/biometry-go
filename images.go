package biometry

import (
	"context"
	"io"
	"net/http"
)

// ImageService handles face image upload and retrieval.
type ImageService struct{ c *Client }

// Upload uploads a face image associated with an email.
func (s *ImageService) Upload(ctx context.Context, email string, image io.Reader, fileName string, opts ...CallOption) (*ImageUploadResponse, error) {
	fields := map[string]string{"email": email}
	files := []FileUpload{{
		FieldName: "image",
		FileName:  fileName,
		Content:   image,
	}}
	var resp ImageUploadResponse
	if err := s.c.doMultipart(ctx, pathGateway+"/images", fields, files, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves the image URL for a given email.
func (s *ImageService) Get(ctx context.Context, email string, opts ...CallOption) (*ImageGetResponse, error) {
	var resp ImageGetResponse
	req, err := s.c.newRequest(ctx, http.MethodGet, pathGateway+"/images/"+email, nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	if err := s.c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
