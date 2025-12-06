package core

type VerifyIdentityRequest struct {
	UserID       int    `json:"user_id"`
	DocumentType string `json:"document_type"` // passport, drivers_license
	DocumentURL  string `json:"document_url"`
}

type IdentityWebhookRequest struct {
	UserID     int    `json:"user_id"`
	DocumentID int    `json:"document_id"`
	Status     string `json:"status"` // passed, failed
}
