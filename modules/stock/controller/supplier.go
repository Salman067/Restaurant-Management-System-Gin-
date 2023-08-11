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

type SupplierControllerInterface interface {
	Pong(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
}

type supplierController struct {
	errors          errors.GinError
	supplierService service.SupplierServiceInterface
}

func NewSupplierController(supplierService service.SupplierServiceInterface) *supplierController {
	return &supplierController{supplierService: supplierService}
}

func (sc *supplierController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from SupplierController")
}
func (sc *supplierController) FindAll(context *gin.Context) {
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
	supplierList, err := sc.supplierService.FindAll(requestParams)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"supplierlist": supplierList})
}

func (sc *supplierController) FindByID(context *gin.Context) {
	supplierID := context.Param("id")

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
	supplier, err := sc.supplierService.FindByID(requestParams, supplierID)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"supplier": supplier})
}
