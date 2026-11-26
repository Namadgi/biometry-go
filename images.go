package biometry

import (
	"context"
	"io"
	"net/http"
)

// ImageService handles a user's reference image (used by face-match).
type ImageService struct{ c *Client }

// Upload stores the reference image for the given user. Requires an approved
// consent for the user.
func (s *ImageService) Upload(ctx context.Context, userID string, image io.Reader, fileName string, opts ...CallOption) error {
	files := []FileUpload{{FieldName: "image", FileName: fileName, Content: image}}
	req, err := s.c.buildMultipartRequestMethod(ctx, http.MethodPut, pathGatewayV2+"/users/"+userID+"/reference-image", nil, files)
	if err != nil {
		return err
	}
	applyCallOptions(req, opts)
	return s.c.do(req, nil)
}

// Get downloads the stored reference image bytes for the given user.
func (s *ImageService) Get(ctx context.Context, userID string, opts ...CallOption) ([]byte, error) {
	req, err := s.c.newRequest(ctx, http.MethodGet, pathGatewayV2+"/users/"+userID+"/reference-image", nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	data, _, err := s.c.doBytes(req)
	return data, err
}
