package biometry

import (
	"context"
	"io"
	"net/http"
)

// FaceMatchService handles face matching between images and videos.
type FaceMatchService struct{ c *Client }

// Match compares faces from an uploaded image and video.
func (f *FaceMatchService) Match(ctx context.Context, image io.Reader, imageName string, video io.Reader, videoName string, opts ...CallOption) (*Response, error) {
	files := []FileUpload{
		{FieldName: "image", FileName: imageName, Content: image},
		{FieldName: "video", FileName: videoName, Content: video},
	}
	return f.c.doMultipartRaw(ctx, pathGateway+"/match-faces", nil, files, opts...)
}

// MatchWithSessionVideo compares a face image against the video from a
// previous process-video call in the given session. This avoids re-uploading
// the video.
func (f *FaceMatchService) MatchWithSessionVideo(ctx context.Context, image io.Reader, imageName string, sessionID string, opts ...CallOption) (*Response, error) {
	files := []FileUpload{
		{FieldName: "image", FileName: imageName, Content: image},
	}
	opts = append(opts, func(r *http.Request) {
		r.Header.Set("X-Use-Prefilled-Video", "true")
		r.Header.Set("X-Session-ID", sessionID)
	})
	return f.c.doMultipartRaw(ctx, pathGateway+"/match-faces", nil, files, opts...)
}
