package schema

import "time"

type Attachment struct {
	ID            uint64 `gorm:"primaryKey"`
	Name          string
	Path          string
	AttachmentKey string
	OwnerID       uint64
	AccountID     uint64
	CreatedAt     time.Time
	CreatedBy     uint64
	UpdatedAt     time.Time
	UpdatedBy     uint64
}
