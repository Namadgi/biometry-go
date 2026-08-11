package biometry

import (
	"encoding/json"
	"fmt"
)

// Excludable verification services for the liveness pipeline.
const (
	ServiceFaceLiveness     = "face_liveness_detection"
	ServiceActiveSpeaker    = "active_speaker_detection"
	ServiceVisualSpeech     = "visual_speech_recognition"
	ServiceFaceRecognition  = "face_recognition"
	ServiceVoiceRecognition = "voice_recognition"
)

// v2Request is the JSON `request` part sent to API v2 gateway endpoints (as a
// multipart form field named "request", or as the JSON body). The gateway
// decodes it and derives the internal identity, session, and service-selection
// headers from it — those headers are stripped from the incoming request, so
// this struct is the only supported way to pass them.
type v2Request struct {
	UserID          string      `json:"user_id,omitempty"`
	SessionID       string      `json:"session_id,omitempty"`
	Phrase          string      `json:"phrase,omitempty"`
	Vocabulary      string      `json:"vocabulary,omitempty"`
	Trigger         string      `json:"trigger,omitempty"`
	UseSessionVideo bool        `json:"use_session_video,omitempty"`
	Provider        string      `json:"provider,omitempty"`
	MRZProvider     string      `json:"mrz_provider,omitempty"`
	IsDocument      bool        `json:"is_document,omitempty"`
	Services        *v2Services `json:"services,omitempty"`
}

type v2Services struct {
	Exclude []string `json:"exclude,omitempty"`
}

// requestField serializes req into the multipart "request" form field.
func requestField(req v2Request) (map[string]string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding request part: %w", err)
	}
	return map[string]string{"request": string(b)}, nil
}
