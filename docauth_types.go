package biometry

// DocAuthInfo contains the extracted document information.
type DocAuthInfo struct {
	DocumentType        string         `json:"document_type"`
	CountryCode         string         `json:"country_code"`
	NationalityCode     string         `json:"nationality_code"`
	NationalityName     string         `json:"nationality_name"`
	Sex                 string         `json:"sex"`
	FirstName           string         `json:"first_name"`
	FatherName          string         `json:"father_name"`
	LastName            string         `json:"last_name"`
	ExpiryDate          string         `json:"expiry_date"`
	DocumentNumber      string         `json:"document_number"`
	BirthDate           string         `json:"birth_date"`
	PortraitPhoto       string         `json:"portrait_photo"`
	Signature           string         `json:"signature"`
	DocumentCategory    string         `json:"document_category"`
	IssuingState        string         `json:"issuing_state"`
	FrontDocumentTypeID string         `json:"front_document_type_id"`
	ContainsRFID        bool           `json:"contains_rfid"`
	FaceImageBase64     string         `json:"face_image_base64,omitempty"`
	CurrentResult       string         `json:"current_result,omitempty"`
	MRZValidation       *MRZValidation `json:"mrz_validation,omitempty"`
}

// MRZValidation contains MRZ check results.
type MRZValidation struct {
	HasMRZ              bool     `json:"has_mrz"`
	RawMRZ              string   `json:"raw_mrz,omitempty"`
	CheckDigitsValid    bool     `json:"check_digits_valid"`
	FieldsMatch         bool     `json:"fields_match"`
	Discrepancies       []string `json:"discrepancies,omitempty"`
	ValidationPerformed bool     `json:"validation_performed"`
}
