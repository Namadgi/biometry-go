package biometry

import "encoding/json"

// Response is the full API response envelope. Used for endpoints like
// process-video and match-faces that return scoring data alongside the result.
type Response struct {
	Data            json.RawMessage `json:"data,omitempty"`
	ScoringResult   json.RawMessage `json:"scoring_result,omitempty"`
	Score           float64         `json:"score,omitempty"`
	DecisionReasons []string        `json:"decision_reasons,omitempty"`
	Message         string          `json:"message,omitempty"`
}

// Decode unmarshals the Data field into the provided target.
func (r *Response) Decode(v any) error {
	return json.Unmarshal(r.Data, v)
}
