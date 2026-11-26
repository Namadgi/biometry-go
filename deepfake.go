package biometry

import (
	"context"
	"io"
	"net/http"
)

// DeepfakeService handles asynchronous deepfake / anti-video-forgery checks.
type DeepfakeService struct{ c *Client }

// DeepfakeInput is the input for submitting a deepfake check.
type DeepfakeInput struct {
	UserID    string // opaque, customer-provided identity key
	Video     io.Reader
	FileName  string
	SessionID string // optional
}

// Submit sends a video for deepfake analysis and returns the created check.
func (d *DeepfakeService) Submit(ctx context.Context, input DeepfakeInput, opts ...CallOption) (*DeepfakeCheck, error) {
	fields, err := requestField(v2Request{UserID: input.UserID, SessionID: input.SessionID})
	if err != nil {
		return nil, err
	}
	files := []FileUpload{{FieldName: "video", FileName: input.FileName, Content: input.Video}}
	var check DeepfakeCheck
	if err := d.c.doMultipart(ctx, pathGatewayV2+"/deepfake/checks", fields, files, &check, opts...); err != nil {
		return nil, err
	}
	return &check, nil
}

// Status retrieves the current state and results of a deepfake check.
func (d *DeepfakeService) Status(ctx context.Context, requestID string, opts ...CallOption) (*DeepfakeCheck, error) {
	req, err := d.c.newRequest(ctx, http.MethodGet, pathGatewayV2+"/deepfake/checks/"+requestID, nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	var check DeepfakeCheck
	if err := d.c.do(req, &check); err != nil {
		return nil, err
	}
	return &check, nil
}
