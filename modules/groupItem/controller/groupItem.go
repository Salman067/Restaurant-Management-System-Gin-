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
	groupItemConsts "pi-inventory/modules/groupItem/consts"
	"pi-inventory/modules/groupItem/models"
	"pi-inventory/modules/groupItem/service"
	"time"

	"github.com/gin-gonic/gin"
)

type GroupItemControllerInterface interface {
	Pong(context *gin.Context)
	CreateGroupItem(context *gin.Context)
	FindByID(context *gin.Context)
	FindAll(context *gin.Context)
	UpdateGroupItem(context *gin.Context)
}

type groupItemController struct {
	page               commonModel.Page
	errors             errors.GinError
	groupItemService   service.GroupItemServiceInterface
	activityLogService commonService.ActivityLogServiceInterface
}

func NewGroupItemController(groupItemService service.GroupItemServiceInterface,
	activityLogService commonService.ActivityLogServiceInterface) *groupItemController {
	return &groupItemController{
		groupItemService:   groupItemService,
		activityLogService: activityLogService,
	}
}

func (gic *groupItemController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from groupItem")
}
func (gic *groupItemController) CreateGroupItem(context *gin.Context) {
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
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
		AccountSlug: accountSlug,
	}

	reqGroupItem := &models.RequestGroupItemBody{}
	if err := context.BindJSON(reqGroupItem); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "group item create"})})
		return
	}
	groupItem, err := gic.groupItemService.CreateGroupItem(requestParams, reqGroupItem, context)
	if err != nil {
		context.AbortWithStatusJSON(gic.errors.GetStatusCode(err), gin.H{"error": gic.errors.ErrorTraverse(err)})
		return
	}

	accountName := context.GetString("account_name")
	notificationText := "A new group item was created "
	activityLogErr := gic.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(groupItem),
		ModelName:        groupItemConsts.ModelNameAddGroupItemRequestBody,
		ModelId:          uint(groupItem.ID),
		ActivityType:     consts.ActionCreate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "group item"}), "id": groupItem.ID})

}

func (gic *groupItemController) FindAll(context *gin.Context) {
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

	queryParamStruct := &models.GroupItemQueryParams{
		KeyWord: context.Query("key_word"),
	}
	page, err := gic.page.GetPageInformation(context)
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

	groupItems, pageInfo, err := gic.groupItemService.FindAll(requestParams, queryParamStruct)
	if err != nil {
		context.AbortWithStatusJSON(gic.errors.GetStatusCode(err), gin.H{"error": gic.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "All groupItems", "groupItemList": groupItems, "Page_info": pageInfo})
}

func (gic *groupItemController) FindByID(context *gin.Context) {
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

	groupItemID, err := utils.Param(context)
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

	groupItem, err := gic.groupItemService.FindByID(requestParams, groupItemID)
	if err != nil {
		context.AbortWithStatusJSON(gic.errors.GetStatusCode(err), gin.H{"error": gic.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"groupItem": groupItem})
}

func (gic *groupItemController) UpdateGroupItem(context *gin.Context) {
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

	groupItemID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		CreatedBy:   uint64(ownerID),
		AccountID:   accountId,
		AccountSlug: accountSlug,
	}
	newGroupItem := &models.UpdateGroupItemRequestBody{}
	if err := context.Bind(&newGroupItem); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "group item update"})})
		return
	}
	updatedGroupItem, err := gic.groupItemService.UpdateGroupItem(requestParams, groupItemID, newGroupItem)
	if err != nil {
		context.AbortWithStatusJSON(gic.errors.GetStatusCode(err), gin.H{"error": gic.errors.ErrorTraverse(err)})
		return
	}
	accountName := context.GetString("account_name")
	notificationText := "A group item was updated"
	activityLogErr := gic.activityLogService.CreateActivityLog(commonModel.ActivityLog{
		CreatedBy:        uint(ownerID),
		AccountId:        accountId,
		NewEntry:         controller.GetAsString(updatedGroupItem),
		ModelName:        groupItemConsts.ModelNameUpdateGroupItemRequestBody,
		ModelId:          uint(updatedGroupItem.ID),
		ActivityType:     consts.ActionUpdate,
		NotificationText: notificationText + accountName,
		SearchBy:         time.Now().String(),
		PushNotification: true,
	})
	if activityLogErr != nil {
		logger.LogError(activityLogErr)
	}
	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "group item"}), "id": updatedGroupItem.ID})
}
