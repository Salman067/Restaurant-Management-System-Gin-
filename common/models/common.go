package models

import "gorm.io/gorm"

type GroupItemAndCompositeItemSummaryResponse struct {
	GroupItemCount     int64 `json:"group_item_count"`
	CompositeItemCount int64 `json:"composite_item_count"`
}

type AccountRequstParams struct {
	CreatedBy   uint64
	AccountSlug string
	AccountID   uint64
	Page        *Page
}

type AccountUserPermission struct {
	gorm.Model
	UserID     uint   `gorm:"index:idx_account_user_permissions_user_id" json:"user_id"`
	AccountID  uint   `gorm:"index:idx_account_user_permissions_account_id" json:"account_id"`
	Permission string `gorm:"size:15"`
}

type User struct {
	gorm.Model
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" gorm:"index:idx_users_email"`
	Mobile         string `json:"mobile" binding:"required" gorm:"unique"`
	Password       string `json:"password"`
	Active         bool   `json:"active"`
	ProfilePicture string `json:"profile_picture"`
	HasBusiness    bool   `json:"has_business"`
}
