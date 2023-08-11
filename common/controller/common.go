package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/logger"
	"pi-inventory/common/models"
	"pi-inventory/common/service"
	"pi-inventory/errors"

	"github.com/gin-gonic/gin"
)

type CommonControllerInterface interface {
	Pong(context *gin.Context)
	GroupItemAndCompositeItemSummary(context *gin.Context)
}

type commonController struct {
	errors        errors.GinError
	commonService service.CommonServiceInterface
}

func NewCommonController(commonService service.CommonServiceInterface) *commonController {
	return &commonController{commonService: commonService}
}

func (cc *commonController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from group and composite summary")
}

func (cc *commonController) GroupItemAndCompositeItemSummary(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo models.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(models.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := models.AccountRequstParams{
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
	}

	groupAndCompositeItemSummary, err := cc.commonService.GroupItemAndCompositeItemSummary(requestParams)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"groupAndCompositeItemSummary": groupAndCompositeItemSummary})
}
