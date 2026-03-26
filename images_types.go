package biometry

// ImageUploadResponse is returned after uploading an image.
type ImageUploadResponse struct {
	URL string `json:"url"`
}

// ImageGetResponse is returned when retrieving an image.
type ImageGetResponse struct {
	URL string `json:"url"`
}
