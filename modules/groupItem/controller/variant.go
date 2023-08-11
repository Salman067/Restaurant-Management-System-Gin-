package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/controller"
	"pi-inventory/common/logger"
	commonService "pi-inventory/common/service"
	"pi-inventory/common/utils"
	"pi-inventory/errors"
	groupItemConsts "pi-inventory/modules/groupItem/consts"
	"pi-inventory/modules/groupItem/models"
	"pi-inventory/modules/groupItem/service"
	"time"

	commonModels "pi-inventory/common/models"

	"github.com/gin-gonic/gin"
)

type VariantControllerInterface interface {
	CreateVariant(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
	DeleteVariant(context *gin.Context)
	UpdateVariant(context *gin.Context)
}

type variantController struct {
	errors             errors.GinError
	variantService     service.VariantServiceInterface
	activityLogService commonService.ActivityLogServiceInterface
}

func NewVariantController(variantSrv service.VariantServiceInterface,
	activityLogService commonService.ActivityLogServiceInterface) *variantController {
	return &variantController{
		variantService:     variantSrv,
		activityLogService: activityLogService,
	}
}

func (vc *variantController) CreateVariant(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModels.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	reqVariant := &models.AddVariantRequestBody{}
	if err := context.BindJSON(&reqVariant); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "variant create"})})
		return
	}
	variant, err := vc.variantService.CreateVariant(requestParams, reqVariant)
	if err != nil {
		context.AbortWithStatusJSON(vc.errors.GetStatusCode(err), gin.H{"error": vc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A new variant was created "
	activityLogErr := vc.activityLogService.CreateActivityLog(commonModels.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(variant),
		ModelName:        groupItemConsts.ModelNameAddVariantRequestBody,
		ModelId:          uint(variant.ID),
		ActivityType:     consts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "attribute"}), "id": variant.ID})
}

func (vc *variantController) FindAll(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModels.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	queryParamStruct := &models.VariantQueryParams{
		KeyWord: context.Query("key_word"),
	}

	variants, err := vc.variantService.FindAll(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(vc.errors.GetStatusCode(err), gin.H{"error": vc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"variantList": variants})
}

func (vc *variantController) FindByID(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
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
	requestParams := commonModels.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	variant, err := vc.variantService.FindByID(requestParams, ID)
	if err != nil {
		context.AbortWithStatusJSON(vc.errors.GetStatusCode(err), gin.H{"error": vc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"singleVariant": variant})
}

func (vc *variantController) DeleteVariant(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
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
	requestParams := commonModels.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	ID, err = vc.variantService.DeleteVariant(requestParams, ID)
	if err != nil {
		context.AbortWithStatusJSON(vc.errors.GetStatusCode(err), gin.H{"error": vc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A variant was deleted "
	activityLogErr := vc.activityLogService.CreateActivityLog(commonModels.ActivityLog{
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
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullDelete", map[string]interface{}{"data": "attribute"}), "id": ID})
}

func (vc *variantController) UpdateVariant(context *gin.Context) {
	//ctx := context.Request.Context()
	accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
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
	requestParams := commonModels.AccountRequstParams{
		AccountSlug: accountSlug,
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
	}
	var variant = new(models.UpdateVariantRequestBody)
	if err := context.Bind(variant); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "variant update"})})
		return
	}
	variants, err := vc.variantService.UpdateVariant(requestParams, ID, variant)
	if err != nil {
		context.AbortWithStatusJSON(vc.errors.GetStatusCode(err), gin.H{"error": vc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A variant was updated"
	activityLogErr := vc.activityLogService.CreateActivityLog(commonModels.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(variants),
		ModelName:        groupItemConsts.ModelNameUpdateVariantRequestBody,
		ModelId:          uint(variants.ID),
		ActivityType:     consts.ActionUpdate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "attribute"}), "id": variants.ID})
}
