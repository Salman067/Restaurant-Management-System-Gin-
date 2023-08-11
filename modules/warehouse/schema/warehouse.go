package schema

import "time"

type Warehouse struct {
	ID        uint64 `gorm:"primaryKey"`
	OwnerID   uint64
	AccountID uint64
	Title     string
	Address   string
	IsUsed    bool
	CreatedAt time.Time
	CreatedBy uint64
	UpdatedAt time.Time
	UpdatedBy uint64
}
