package controller

import (
	"net/http"
	commonConsts "pi-inventory/common/consts"
	"pi-inventory/common/controller"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	commonService "pi-inventory/common/service"
	"pi-inventory/common/utils"
	"pi-inventory/errors"
	"pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StockActivityControllerInterface interface {
	Pong(context *gin.Context)
	CreateStockActivity(context *gin.Context)
	FindAll(context *gin.Context)
	CreateBulkStockActivity(context *gin.Context)
}

type stockActivityController struct {
	page                 commonModel.Page
	errors               errors.GinError
	stockActivityService service.StockActivityServiceInterface
	activityLogService   commonService.ActivityLogServiceInterface
}

func NewStockActivityController(service service.StockActivityServiceInterface,
	activityLogService commonService.ActivityLogServiceInterface) *stockActivityController {
	return &stockActivityController{
		stockActivityService: service,
		activityLogService:   activityLogService,
	}
}

func (sac *stockActivityController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from stockActivityController")
}

func (sac *stockActivityController) CreateStockActivity(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if commonConsts.AccountTypes[accountInfo.Type] == commonConsts.AccountTypeBusinessSubBook || commonConsts.AccountTypes[accountInfo.Type] == commonConsts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}

	reqStockActivity := &models.AddStockActivityRequestBody{}

	err := context.BindJSON(&reqStockActivity)
	if err != nil && reqStockActivity.Mode == consts.StockActivityModeQuantity {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError1", nil)})
		return
	}
	stockActivity, err := sac.stockActivityService.CreateStockActivity(requestParams, reqStockActivity)
	if err != nil {
		context.AbortWithStatusJSON(sac.errors.GetStatusCode(err), gin.H{"error": sac.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A new stock activity was created "
	activityLogErr := sac.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(stockActivity),
		ModelName:        consts.ModelNameAddStockActivityRequestBody,
		ModelId:          uint(stockActivity.ID),
		ActivityType:     commonConsts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "stock activity"}), "id": stockActivity.ID})
}

func (sac *stockActivityController) FindAll(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if commonConsts.AccountTypes[accountInfo.Type] == commonConsts.AccountTypeBusinessSubBook || commonConsts.AccountTypes[accountInfo.Type] == commonConsts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	stringStockID := context.Query("stock")
	stockID64, err := strconv.ParseInt(stringStockID, 10, 64)
	if err != nil && stringStockID != "" {
		logger.LogError("Stock id invalid")
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidQueryParam", nil)})
		return
	}
	var queryParamStruct models.StockActivityQueryParams
	if stringStockID == "" {
		queryParamStruct = models.StockActivityQueryParams{
			StockID: 0,
			Keyword: context.Query("key_word"),
		}
	} else {
		queryParamStruct = models.StockActivityQueryParams{
			StockID: uint64(stockID64),
			Keyword: context.Query("key_word"),
		}
	}

	page, err := sac.page.GetPageInformation(context)
	if err != nil {
		logger.LogError("Invalid page information", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidPageInformation", nil)})
		return
	}
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
		Page:        page,
	}
	stocks, pageInfo, err := sac.stockActivityService.FindAll(requestParams, &queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(sac.errors.GetStatusCode(err), gin.H{"error": sac.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "All stocks", "stockList": stocks, "Page_info": pageInfo})

}

func (sac *stockActivityController) CreateBulkStockActivity(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	logger.LogError(accountSlug)
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if commonConsts.AccountTypes[accountInfo.Type] == commonConsts.AccountTypeBusinessSubBook || commonConsts.AccountTypes[accountInfo.Type] == commonConsts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}

	reqBulkStockActivity := &models.BulkAdjustment{}
	if err := context.BindJSON(&reqBulkStockActivity); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError1", nil)})
		return
	}
	_, err := sac.stockActivityService.CreateBulkStockActivity(requestParams, reqBulkStockActivity)
	if err != nil {
		context.AbortWithStatusJSON(sac.errors.GetStatusCode(err), gin.H{"error": sac.errors.ErrorTraverse(err)})
		return
	}
	//context.JSON(http.StatusCreated, gin.H{"message": "CreatedSuccessfully", "id": stockActivity.ID})
}
