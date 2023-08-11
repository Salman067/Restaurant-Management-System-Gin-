package middlewares

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/logger"
	"pi-inventory/common/models"
	commonModuleService "pi-inventory/common/service"
	"pi-inventory/common/utils"
	"pi-inventory/dic"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AccountValidate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountId := ctx.Params.ByName("id")
		if accountId == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": utils.Trans("invalidAccountID", nil)})
			return
		}

		// Validate Account permission
		ValidateUserPermission(ctx, accountId, "id")
		ctx.Next()
	}
}

func AccountDependValidate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountSlug := ctx.Params.ByName("account_slug")

		if accountSlug == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": utils.Trans("invalidAccountID", nil)})
			return
		}

		// Validate Account permission
		ValidateUserPermission(ctx, accountSlug, "slug")

		ctx.Next()
	}
}

func ValidateBusinessPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if consts.AccountTypes[ctx.GetString("account_type")] == consts.AccountTypePersonal {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": utils.Trans("notBusinessBook", nil)})
			return
		}

		ctx.Next()
	}
}

func ValidateUserPermission(ctx *gin.Context, accountId string, findBy string) {

	repo := dic.Container.Get(consts.RedisCommonService).(*commonModuleService.RedisService)

	accountInfo := models.RedisAccountInfo{}
	//repo.GetAll(ctx, "product_")
	if findBy == "slug" {
		_, err := repo.GetRedisAccountInfo(&accountInfo, accountId)
		if err != nil {
			logger.LogError(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
	} else {
		accId, _ := strconv.Atoi(accountId)
		accountInfo.ID = uint(accId)
		_, err := repo.GetRedisAccountInfo(&accountInfo, "")
		if err != nil {
			logger.LogError(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	if accountInfo.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": utils.Trans("invalidAccountID", nil)})
		return
	}

	if !accountInfo.Enabled {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "error", "message": utils.Trans("cashbookIsDisabled", nil)})
		return
	}

	// permission, err := accountInfo.GetUserPermission(uint(ctx.GetInt64("user_id")))
	// if err != nil {
	// 	logger.LogError(err)
	// 	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
	// 	return
	// }
	logger.LogError(accountInfo)

	ctx.Set("account_info", accountInfo)
	ctx.Set("account_id", uint64(accountInfo.ID))
	//ctx.Set("user_access", permission)
	ctx.Set("account_name", accountInfo.Name)
	ctx.Set("account_created_by", accountInfo.CreatedBy)
	ctx.Set("rows_per_page", accountInfo.RowsPerPage)
	ctx.Set("last_voucher_number", accountInfo.LastVoucherNumber)
	ctx.Set("voucher_prefix", accountInfo.VoucherPrefix)
	ctx.Set("allow_negative_balance", accountInfo.AllowNegativeBalance)
	ctx.Set("enable_cash_in_notification_for_amount_greater_than", accountInfo.EnableCashInNotificationForAmountGreaterThan)
	ctx.Set("enable_cash_out_notification", accountInfo.EnableCashOutNotification)
	ctx.Set("account_type", accountInfo.Type)
	ctx.Set("main_account_id", accountInfo.MainAccountId)
	ctx.Set("branch_account_id", accountInfo.BranchAccountId)
	ctx.Set("enable_cash_out_approval", accountInfo.EnableCashOutApproval)
	ctx.Set("cash_out_approval_amount", accountInfo.CashOutApprovalAmount)
	if accountInfo.BusinessId != nil {
		ctx.Set("business_id", uint64(*accountInfo.BusinessId))
	}

	return
}

func AccountAdminPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission := ctx.GetString("user_access")
		if permission != consts.AccountPermissionAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": utils.Trans("userDoesNotHavePermissionToAccess", nil)})
			return
		}
		ctx.Next()
	}
}

func AccountOperatorPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission := ctx.GetString("user_access")
		switch permission {
		case
			consts.AccountPermissionAdmin,
			consts.AccountPermissionApprover,
			consts.AccountPermissionOperator:
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": utils.Trans("userDoesNotHavePermissionToAccess", nil)})
		return
	}
}

func AccountApproverPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission := ctx.GetString("user_access")
		switch permission {
		case
			consts.AccountPermissionAdmin,
			consts.AccountPermissionApprover:
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": utils.Trans("userDoesNotHavePermissionToAccess", nil)})
		return
	}
}

func BusinessAccountPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountType := ctx.GetString("account_type")
		switch accountType {
		case
			consts.AccountTypesToString[consts.AccountTypeBusinessMainBook],
			consts.AccountTypesToString[consts.AccountTypeBusinessBranchBook],
			consts.AccountTypesToString[consts.AccountTypeBusinessSubBook]:
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": utils.Trans("accountDoesNotHavePermissionToAccess", nil)})
		return
	}
}
