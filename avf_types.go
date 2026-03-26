package biometry

// AVFInferenceResponse is returned when starting an AVF inference.
type AVFInferenceResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

// AVFStatusResponse is returned when checking AVF inference status.
type AVFStatusResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Result    any    `json:"result,omitempty"`
}
