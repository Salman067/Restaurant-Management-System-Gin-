package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/controller"
	"pi-inventory/common/logger"
	commonService "pi-inventory/common/service"
	"pi-inventory/common/utils"
	"pi-inventory/errors"
	warehouseConsts "pi-inventory/modules/warehouse/consts"
	"pi-inventory/modules/warehouse/models"
	"pi-inventory/modules/warehouse/service"
	"time"

	commonModel "pi-inventory/common/models"

	"github.com/gin-gonic/gin"
)

type WarehouseControllerInterface interface {
	Pong(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
	CreateWarehouse(context *gin.Context)
	UpdateWarehouse(context *gin.Context)
	DeleteWarehouse(context *gin.Context)
}

type warehouseController struct {
	page               commonModel.Page
	errors             errors.GinError
	warehouseService   service.WarehouseServiceInterface
	activityLogService commonService.ActivityLogServiceInterface
}

func NewWarehouseController(warehouseSrv service.WarehouseServiceInterface,
	activityLogService commonService.ActivityLogServiceInterface) *warehouseController {
	return &warehouseController{
		warehouseService:   warehouseSrv,
		activityLogService: activityLogService,
	}
}

func (cc *warehouseController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from Warehouse")
}

func (cc *warehouseController) FindAll(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	queryParamStruct := &models.WarehouseQueryParams{
		KeyWord: context.Query("key_word"),
	}
	page, err := cc.page.GetPageInformation(context)
	if err != nil {
		logger.LogError("Invalid page information", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidPageInformation", nil)})
		return
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
		Page:        page,
	}
	warehouses, pageInfo, err := cc.warehouseService.FindAll(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "All warehouses", "warehouseList": warehouses, "Page_info": pageInfo})
}

func (cc *warehouseController) FindByID(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	warehouseID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	warehouse, err := cc.warehouseService.FindByID(requestParams, warehouseID)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"warehouse": warehouse})
}

func (cc *warehouseController) CreateWarehouse(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	reqWarehouse := &models.AddWarehouseRequestBody{}
	if err := context.BindJSON(&reqWarehouse); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "location create"})})
		return
	}
	warehouse, err := cc.warehouseService.CreateWarehouse(requestParams, reqWarehouse)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A new warehouse was created "
	activityLogErr := cc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(warehouse),
		ModelName:        warehouseConsts.ModelNameAddWarehouseRequestBody,
		ModelId:          uint(warehouse.ID),
		ActivityType:     consts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "location"}), "id": warehouse.ID})
}

func (cc *warehouseController) UpdateWarehouse(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	warehouseID, err := utils.Param(context)
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

	reqBody := &models.UpdateWarehouseRequestBody{}
	if err := context.BindJSON(&reqBody); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "location update"})})
		return
	}
	warehouse, err := cc.warehouseService.UpdateWarehouse(requestParams, warehouseID, reqBody)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A warehouse was updated "
	activityLogErr := cc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(warehouse),
		ModelName:        warehouseConsts.ModelNameUpdateWarehouseRequestBody,
		ModelId:          uint(warehouse.ID),
		ActivityType:     consts.ActionUpdate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "location"}), "id": warehouse.ID})
}

func (cc *warehouseController) DeleteWarehouse(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	warehouseID, err := utils.Param(context)
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

	ID, err := cc.warehouseService.DeleteWarehouse(requestParams, warehouseID)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A warehouse was deleted "
	activityLogErr := cc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		ModelId:          uint(ID),
		ActivityType:     consts.ActionDelete,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullDelete", map[string]interface{}{"data": "location"}), "id": ID})
}
