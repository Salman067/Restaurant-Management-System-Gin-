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
	stockModel "pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/service"
	"time"

	"github.com/gin-gonic/gin"
)

type StockControllerInterface interface {
	Pong(context *gin.Context)
	CreateStock(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
	UpdateStock(context *gin.Context)
	DeleteStock(context *gin.Context)
	StockSummary(context *gin.Context)
}

type StockController struct {
	page               commonModel.Page
	errors             errors.GinError
	stockService       service.StockServiceInterface
	activityLogService commonService.ActivityLogServiceInterface
}

func NewStockController(service service.StockServiceInterface, activityLogService commonService.ActivityLogServiceInterface) *StockController {
	return &StockController{stockService: service, activityLogService: activityLogService}
}

func (sc *StockController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from stock")
}

func (sc *StockController) CreateStock(context *gin.Context) {
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

	reqStock := &stockModel.AddStockRequestBody{}
	if err := context.BindJSON(&reqStock); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "stock create"})})
		return
	}

	stock, err := sc.stockService.CreateStock(requestParams, reqStock, context)
	if err != nil {
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A new stock was created "
	activityLogErr := sc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(stock),
		ModelName:        consts.ModelNameAddStockRequestBody,
		ModelId:          uint(stock.ID),
		ActivityType:     commonConsts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "item"}), "id": stock.ID})
}

func (sc *StockController) FindAll(context *gin.Context) {
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

	queryParamStruct := &stockModel.StockQueryParams{
		Type:     context.Query("type"),
		Status:   context.Query("status"),
		ListType: context.Query("list_type"),
		KeyWord:  context.Query("key_word"),
	}
	page, err := sc.page.GetPageInformation(context)
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

	stocks, pageInfo, err := sc.stockService.FindAllNew(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "All stocks", "stockList": stocks, "Page_info": pageInfo})
}

func (sc *StockController) FindByID(context *gin.Context) {
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

	stockID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	stock, err := sc.stockService.FindByID(requestParams, stockID)
	if err != nil {
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Stock": stock})
}

func (sc *StockController) UpdateStock(context *gin.Context) {

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

	ID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	var stock = new(stockModel.UpdateStockRequestBody)
	if err := context.Bind(stock); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "stock update"})})
		return
	}
	stocks, err := sc.stockService.UpdateStock(requestParams, ID, stock)
	if err != nil {
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A stock was updated "
	activityLogErr := sc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(stocks),
		ModelName:        consts.ModelNameUpdateStockRequestBody,
		ModelId:          uint(stocks.ID),
		ActivityType:     commonConsts.ActionUpdate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "item"}), "id": stocks.ID})
}

func (sc *StockController) StockSummary(context *gin.Context) {
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

	page, err := sc.page.GetPageInformation(context)
	if err != nil {
		logger.LogError("Invalid page information", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidPageInformation", nil)})
		return
	}
	queryParamStruct := &stockModel.StockQueryParams{
		KeyWord: context.Query("key_word"),
	}

	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
		Page:        page,
	}
	stockSummary, err := sc.stockService.StockSummary(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"stockSummary": stockSummary})
}

func (sc *StockController) DeleteStock(context *gin.Context) {
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

	stockID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}

	stock, err := sc.stockService.DeleteStock(requestParams, stockID)
	if err != nil {
		context.AbortWithStatusJSON(sc.errors.GetStatusCode(err), gin.H{"error": sc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A stock was deleted "
	activityLogErr := sc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		ModelId:          uint(stock.ID),
		ActivityType:     commonConsts.ActionDelete,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullDelete", map[string]interface{}{"data": "item"}), "id": stock.ID})
}
