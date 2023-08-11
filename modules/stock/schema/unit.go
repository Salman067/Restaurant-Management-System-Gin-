package schema

import "time"

type Unit struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Title     string
	Status    string `gorm:"default:active"`
	OwnerID   uint64
	AccountID uint64
	IsUsed    bool
	CreatedAt time.Time
	CreatedBy uint64
	UpdatedAt time.Time
	UpdatedBy uint64
}
