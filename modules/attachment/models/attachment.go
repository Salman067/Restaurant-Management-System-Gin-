package models

type AttachmentCustomResponse struct {
	ID   uint64 `json:"id"`
	Path string `json:"path"`
}
type UploadAttachmentRequestBody struct {
	AttachmentType string `json:"attachment_type"`
	Path           string `json:"path"`
	Name           string `json:"name"`
}

type RequstParams struct {
	CreatedBy   uint64
	AccountSlug string
	AccountID   uint64
}
