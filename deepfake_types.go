package biometry

// DeepfakeCheck is the state of an asynchronous deepfake check.
type DeepfakeCheck struct {
	ID     string `json:"id"`
	Status string `json:"status"` // pending | processing | completed | failed
	// Score is the deepfake probability (0.0–1.0), nil until completed.
	Score       *float64 `json:"score,omitempty"`
	Prediction  string   `json:"prediction,omitempty"`
	Message     string   `json:"message,omitempty"`
	CreatedAt   int64    `json:"created_at,omitempty"`   // unix seconds
	CompletedAt int64    `json:"completed_at,omitempty"` // unix seconds
}
