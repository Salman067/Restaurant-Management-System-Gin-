package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"

	"pi-inventory/common/utils"
	"pi-inventory/errors"
	"pi-inventory/modules/composite/models"
	"pi-inventory/modules/composite/service"

	"github.com/gin-gonic/gin"
)

type CompositeControllerInterface interface {
	Pong(context *gin.Context)
	FindAll(context *gin.Context)
	FindByID(context *gin.Context)
	CreateComposite(context *gin.Context)
	UpdateComposite(context *gin.Context)
}

type compositeController struct {
	page             commonModel.Page
	errors           errors.GinError
	compositeService service.CompositeServiceInterface
}

func NewCompositeController(CompositeSrv service.CompositeServiceInterface) *compositeController {
	return &compositeController{compositeService: CompositeSrv}
}

func (cc *compositeController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from Composite")
}
func (cc *compositeController) FindAll(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	queryParamStruct := &models.CompositeQueryParams{
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
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
		Page:      page,
	}
	Composites, pageInfo, err := cc.compositeService.FindAll(requestParams, queryParamStruct)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"_Message": "All Composites", "compositeList": Composites, "_Page_info": pageInfo})
}

func (cc *compositeController) FindByID(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
	}

	compositeID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	composite, err := cc.compositeService.FindByID(requestParams, compositeID)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Composite": composite})
}

func (cc *compositeController) CreateComposite(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
	}

	reqComposite := &models.AddCompositeRequestBody{}
	if err := context.BindJSON(reqComposite); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "composite create"})})
		return
	}
	composite, err := cc.compositeService.CreateComposite(requestParams, reqComposite)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": utils.Trans("successfullCreate", map[string]interface{}{"data": "composite item"}), "id": composite.ID})
}

func (cc *compositeController) UpdateComposite(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	compositeID, err := utils.Param(context)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("InvalidParams", nil)})
		return
	}
	ownerID := context.GetInt64("user_id")
	requestParams := commonModel.AccountRequstParams{
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
	}

	var reqComposite = new(models.UpdateCompositeRequestBody)
	if err := context.Bind(reqComposite); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", map[string]interface{}{"data": "composite update"})})
		return
	}
	composite, err := cc.compositeService.UpdateComposite(requestParams, compositeID, reqComposite)
	if err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(cc.errors.GetStatusCode(err), gin.H{"error": cc.errors.ErrorTraverse(err)})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"message": utils.Trans("successfullUpdate", map[string]interface{}{"data": "composite item"}), "id": composite.ID})
}
