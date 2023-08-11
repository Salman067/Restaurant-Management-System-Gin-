package models

import (
	"errors"
	"pi-inventory/common/logger"
	"pi-inventory/common/utils"
	"time"
)

type RedisAccountInfo struct {
	ID                                           uint
	Slug                                         string  `json:"slug"`
	Name                                         string  `json:"name"`
	Logo                                         *string `json:"logo"`
	Address                                      string  `json:"address"`
	Phone                                        string  `json:"phone"`
	Type                                         string  `json:"type"`
	MainAccountId                                uint    `json:"main_account_id"`
	MainAccountName                              string  `json:"main_account_name"`
	MainAccountSlug                              string  `json:"main_account_slug"`
	BranchAccountId                              uint    `json:"branch_account_id"`
	BranchAccountName                            string  `json:"branch_account_name"`
	CreatedBy                                    uint64  `json:"created_by"`
	CreatedByName                                string  `json:"created_by_name"`
	RowsPerPage                                  int     `json:"rows_per_page"`
	AllowNegativeBalance                         bool    `json:"allow_negative_balance"`
	EnableCashInNotificationForAmountGreaterThan float64 `json:"enable_cash_in_notification_if_amount_greater_than"`
	EnableCashOutNotification                    bool    `json:"enable_cash_out_notification"  gorm:"default:false"`
	LastVoucherNumber                            uint64  `json:"last_voucher_number"`
	VoucherPrefix                                string  `json:"voucher_prefix"`
	LastEstimateNumber                           uint64  `json:"last_estimate_number"`
	EstimatePrefix                               string  `json:"estimate_prefix"`
	LastInvoiceNumber                            uint64  `json:"last_invoice_number"`
	InvoicePrefix                                string  `json:"invoice_prefix"`
	LastPaymentNumber                            uint64  `json:"last_payment_number"`
	PaymentPrefix                                string  `json:"payment_prefix"`
	EnableCashOutApproval                        bool    `json:"enable_cash_out_approval" gorm:"default:false"`
	CashOutApprovalAmount                        float64 `json:"cash_out_approval_amount"`
	Enabled                                      bool    `json:"enabled"`
	BusinessId                                   *uint   `json:"business_id"`
	AccountUserPermissions                       []RedisAccountUserPermission
}

type RedisAccountUserPermission struct {
	ID         uint
	UserId     uint       `json:"user_id"`
	AccountId  uint       `json:"account_id"`
	DeletedAt  *time.Time `json:"-"`
	Permission string     `json:"permission"`
	User       RedisUser
}

type RedisUser struct {
	ID             uint
	Name           string `json:"name"`
	Email          string `json:"email"`
	Mobile         string `json:"mobile"`
	ProfilePicture string `json:"profile_picture"`
}

func (r *RedisAccountInfo) GetUserPermission(userId uint) (string, error) {
	userPermissions := r.AccountUserPermissions
	for _, userPermission := range userPermissions {
		if userPermission.UserId == userId {
			return userPermission.Permission, nil
		}
	}
	return "", errors.New(utils.Trans("userDoesNotHavePermissionToAccess", nil))
}

func (r *RedisAccountInfo) UpdatePermission(accountInfo *RedisAccountInfo) error {
	accountPermissions := []RedisAccountUserPermission{}
	for _, permission := range accountInfo.AccountUserPermissions {
		redisPermission := RedisAccountUserPermission{}
		redisPermission.ID = permission.ID
		redisPermission.Permission = permission.Permission

		redisPermission.UserId = permission.UserId
		redisPermission.AccountId = permission.AccountId
		redisPermission.User = permission.User
		logger.LogError(redisPermission.User)
		accountPermissions = append(accountPermissions, redisPermission)
	}
	r.AccountUserPermissions = accountPermissions
	return nil
}
