package biometry

import "time"

// TransactionListRequest specifies filters for listing transactions.
type TransactionListRequest struct {
	ProjectID    string `json:"project_id,omitempty"`
	SessionID    string `json:"session_id,omitempty"`
	UserFullname string `json:"user_fullname,omitempty"`
	FromDate     string `json:"from_date,omitempty"`
	ToDate       string `json:"to_date,omitempty"`
	StatusCode   *int   `json:"status_code,omitempty"`
	Decision     string `json:"decision,omitempty"`
	Action       string `json:"action,omitempty"`
	Flag         string `json:"flag,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	SortBy       string `json:"sort_by,omitempty"`
	SortOrder    string `json:"sort_order,omitempty"`
}

// Transaction represents a biometric transaction record.
type Transaction struct {
	ID              string            `json:"id"`
	SessionID       string            `json:"session_id"`
	Flag            string            `json:"flag"`
	Device          TransactionDevice `json:"device"`
	ProjectID       string            `json:"project_id"`
	ApiKeyID        string            `json:"api_key_id"`
	StatusCode      int               `json:"status_code"`
	IsSuccess       bool              `json:"is_success"`
	BodyByteSize    int64             `json:"body_byte_size"`
	UsedServices    []string          `json:"used_services"`
	ReceivedAt      time.Time         `json:"recieved_at"` // note: backend spelling
	AnsweredAt      time.Time         `json:"answered_at"`
	Action          string            `json:"action"`
	AppVersion      string            `json:"app_version"`
	Duration        int64             `json:"duration"`
	Decision        string            `json:"decision"`
	DecisionHistory []DecisionChange  `json:"decision_history"`
	Geo             *Geo              `json:"geo"`
	Tags            []string          `json:"tags"`
	InputMetadata   any               `json:"input_metadata"`
	OutputMetadata  any               `json:"output_metadata"`
	RulesVersion    int               `json:"rules_version"`
}

// TransactionDevice contains device information from a transaction.
type TransactionDevice struct {
	DeviceType string `json:"device_type"`
	DeviceOS   string `json:"device_os"`
}

// DecisionChange records a change in transaction decision.
type DecisionChange struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Reason    string    `json:"reason"`
	ChangedBy string    `json:"changed_by"`
	ChangedAt time.Time `json:"changed_at"`
}

// Geo contains geographic location data.
type Geo struct {
	IP      string  `json:"ip"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Country string  `json:"country"`
	City    string  `json:"city"`
}
