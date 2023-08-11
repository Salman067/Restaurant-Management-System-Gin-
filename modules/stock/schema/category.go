package schema

import "time"

type Category struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string
	OwnerID     uint64
	Status      string `gorm:"default:active"`
	Description string
	IsUsed      bool
	CreatedAt   time.Time
	CreatedBy   uint64
	UpdatedAt   time.Time
	UpdatedBy   uint64
	AccountID   uint64
}
