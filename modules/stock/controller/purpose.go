package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/controller"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	commonService "pi-inventory/common/service"
	"pi-inventory/common/utils"
	"pi-inventory/errors"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/service"
	"time"

	"github.com/gin-gonic/gin"
)

type PurposeControllerInterface interface {
	Pong(context *gin.Context)
	CreatePurpose(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
	DeletePurpose(context *gin.Context)
	UpdatePurpose(context *gin.Context)
}

type purposeController struct {
	errors             errors.GinError
	purposeService     service.PurposeServiceInterface
	activityLogService commonService.ActivityLogServiceInterface
}

func NewPurposeController(purposeService service.PurposeServiceInterface,
	activityLogService commonService.ActivityLogServiceInterface) *purposeController {
	return &purposeController{
		purposeService:     purposeService,
		activityLogService: activityLogService,
	}
}

func (pc *purposeController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from stockActivityController")
}

func (pc *purposeController) CreatePurpose(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	logger.LogError(accountSlug)
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
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
		AccountSlug: accountSlug,
	}

	reqPurpose := &models.AddPurposeRequestBody{}
	if err := context.BindJSON(&reqPurpose); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "reason create"})})
		return
	}
	purpose, err := pc.purposeService.CreatePurpose(requestParams, reqPurpose)
	if err != nil {
		context.AbortWithStatusJSON(pc.errors.GetStatusCode(err), gin.H{"error": pc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A new reason was created "
	activityLogErr := pc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(purpose),
		ModelName:        stockConst.ModelNameAddPurposeRequestBody,
		ModelId:          uint(purpose.ID),
		ActivityType:     consts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "reason"}), "id": purpose.ID})
}

func (pc *purposeController) FindAll(context *gin.Context) {
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
	queryParamStruct := &models.PurposeQueryParams{
		KeyWord: context.Query("key_word"),
	}

	purposes, err := pc.purposeService.FindAll(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(pc.errors.GetStatusCode(err), gin.H{"error": pc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"purposeList": purposes})
}

func (pc *purposeController) FindByID(context *gin.Context) {
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

	purpose, err := pc.purposeService.FindByID(requestParams, ID)
	if err != nil {
		context.AbortWithStatusJSON(pc.errors.GetStatusCode(err), gin.H{"error": pc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"singlePurpose": purpose})
}

func (pc *purposeController) DeletePurpose(context *gin.Context) {
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

	ID, err = pc.purposeService.DeletePurpose(requestParams, ID)
	if err != nil {
		context.AbortWithStatusJSON(pc.errors.GetStatusCode(err), gin.H{"error": pc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A reason was deleted "
	activityLogErr := pc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
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

	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullDelete", map[string]interface{}{"data": "reason"}), "id": ID})
}

func (pc *purposeController) UpdatePurpose(context *gin.Context) {
	ID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	logger.LogError(accountSlug)
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
	var purpose = new(models.UpdatePurposeRequestBody)
	if err := context.Bind(purpose); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "reason update"})})
		return
	}
	purposes, err := pc.purposeService.UpdatePurpose(requestParams, ID, purpose)
	if err != nil {
		context.AbortWithStatusJSON(pc.errors.GetStatusCode(err), gin.H{"error": pc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A reason was updated "
	activityLogErr := pc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(purposes),
		ModelName:        stockConst.ModelNameUpdatePurposeRequestBody,
		ModelId:          uint(purposes.ID),
		ActivityType:     consts.ActionUpdate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "reason"}), "id": purposes.ID})
}
