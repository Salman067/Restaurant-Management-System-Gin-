package schema

import "time"

type Variant struct {
	ID        uint64 `gorm:"primaryKey"`
	Title     string
	OwnerID   uint64
	AccountID uint64
	CreatedAt time.Time
	CreatedBy uint64
	UpdatedAt time.Time
	UpdatedBy uint64
}
