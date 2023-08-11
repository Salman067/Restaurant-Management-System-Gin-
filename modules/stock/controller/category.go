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
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/service"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryControllerInterface interface {
	Pong(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
	CreateCategory(context *gin.Context)
	UpdateCategory(context *gin.Context)
	DeleteCategory(context *gin.Context)
}

type categoryController struct {
	page               commonModel.Page
	errors             errors.GinError
	categoryService    service.CategoryServiceInterface
	activityLogService commonService.ActivityLogServiceInterface
}

func NewCategoryController(categorySrv service.CategoryServiceInterface,
	activityLogService commonService.ActivityLogServiceInterface) *categoryController {
	return &categoryController{
		categoryService:    categorySrv,
		activityLogService: activityLogService,
	}
}

func (cc *categoryController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from category")
}

func (cc *categoryController) FindAll(context *gin.Context) {
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

	queryParamStruct := &models.CategoryQueryParams{
		KeyWord: context.Query("key_word"),
	}
	page, err := cc.page.GetPageInformation(context)
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

	categories, pageInfo, err := cc.categoryService.FindAll(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "All categories", "categories": categories, "Page_info": pageInfo})
}

func (cc *categoryController) FindByID(context *gin.Context) {
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

	categoryID, err := utils.Param(context)
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

	category, err := cc.categoryService.FindByID(requestParams, categoryID)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"category": category})
}

func (cc *categoryController) CreateCategory(context *gin.Context) {
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

	reqCategory := &models.AddCategoryRequestBody{}
	if err := context.BindJSON(&reqCategory); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "category create"})})
		return
	}
	category, err := cc.categoryService.CreateCategory(requestParams, reqCategory)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A new category was created "
	activityLogErr := cc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(category),
		ModelName:        stockConst.ModelNameAddCategoryRequestBody,
		ModelId:          uint(category.ID),
		ActivityType:     commonConsts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "category"}), "id": category.ID})
}

func (cc *categoryController) UpdateCategory(context *gin.Context) {
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

	categoryID, err := utils.Param(context)
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

	reqBody := &models.UpdateCategoryRequestBody{}
	if err := context.BindJSON(&reqBody); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "category update"})})
		return
	}
	category, err := cc.categoryService.UpdateCategory(requestParams, categoryID, reqBody)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A new category was updated "
	activityLogErr := cc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(category),
		ModelName:        stockConst.ModelNameUpdateCategoryRequestBody,
		ModelId:          uint(category.ID),
		ActivityType:     commonConsts.ActionUpdate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "category"}), "id": category.ID})
}

func (cc *categoryController) DeleteCategory(context *gin.Context) {
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

	categoryID, err := utils.Param(context)
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

	categoryID, err = cc.categoryService.DeleteCategory(requestParams, categoryID)
	if err != nil {
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A category was deleted "
	activityLogErr := cc.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		ModelId:          uint(categoryID),
		ActivityType:     commonConsts.ActionDelete,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}

	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullDelete", map[string]interface{}{"data": "category"}), "id": categoryID})
}
