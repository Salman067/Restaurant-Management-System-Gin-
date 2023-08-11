package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	commonModels "pi-inventory/common/models"

	"pi-inventory/common/logger"
	"pi-inventory/errors"
	stockModel "pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/service"

	"github.com/gin-gonic/gin"
)

type TaxControllerInterface interface {
	Pong(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
}

type taxController struct {
	errors     errors.GinError
	taxService service.TaxServiceInterface
}

func NewTaxController(taxService service.TaxServiceInterface) *taxController {
	return &taxController{taxService: taxService}
}

func (tc *taxController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from taxController")
}
func (tc *taxController) FindAll(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	logger.LogError(accountSlug)
	//accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		//accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := stockModel.RequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
	}

	taxList, err := tc.taxService.FindAll(requestParams)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(tc.errors.GetStatusCode(err), gin.H{"error": tc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"taxlist": taxList})
}

func (tc *taxController) FindByID(context *gin.Context) {
	taxID := context.Param("id")
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	logger.LogError(accountSlug)
	//accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		//accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := stockModel.RequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
	}
	tax, err := tc.taxService.FindByID(requestParams, taxID)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(tc.errors.GetStatusCode(err), gin.H{"error": tc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"tax": tax})
}
