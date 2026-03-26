package biometry

import (
	"context"
	"io"
	"net/http"
)

// AVFService handles anti-video-forgery inference.
type AVFService struct{ c *Client }

// StartInference submits a video for anti-forgery analysis.
func (a *AVFService) StartInference(ctx context.Context, video io.Reader, videoName string, opts ...CallOption) (*AVFInferenceResponse, error) {
	files := []FileUpload{{
		FieldName: "video",
		FileName:  videoName,
		Content:   video,
	}}
	var resp AVFInferenceResponse
	if err := a.c.doMultipart(ctx, pathGateway+"/avf/infer", nil, files, &resp, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CheckStatus checks the status of an AVF inference request.
func (a *AVFService) CheckStatus(ctx context.Context, requestID string, opts ...CallOption) (*AVFStatusResponse, error) {
	var resp AVFStatusResponse
	req, err := a.c.newRequest(ctx, http.MethodGet, pathGateway+"/avf/status/"+requestID, nil)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, opts)
	if err := a.c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
