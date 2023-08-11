package controller

import (
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/logger"
	"pi-inventory/common/utils"
	"pi-inventory/modules/attachment/models"
	"pi-inventory/modules/attachment/service"
	"strings"

	commonModels "pi-inventory/common/models"

	"github.com/gin-gonic/gin"
)

type AttachmentControllerInterface interface {
	GetSingleFile(context *gin.Context)
	UploadSingleFile(context *gin.Context)
	DeleteSingleFile(context *gin.Context)
	GetSingleAttachmentFile(context *gin.Context)
	FetchAttachments(context *gin.Context)
	UploadMultipleAttachment(context *gin.Context)
}

type attachmentController struct {
	attachmentService service.AttachmentServiceInterface
}

func NewAttachmentController(service service.AttachmentServiceInterface) *attachmentController {
	return &attachmentController{attachmentService: service}
}

func (ac *attachmentController) GetSingleFile(context *gin.Context) {
	path := context.Query("attachment_path")
	ac.attachmentService.GetSingleAttachment(context, path)
}

func (ac *attachmentController) DeleteSingleFile(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	//accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		//accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	path := context.Query("attachment_path")
	err := ac.attachmentService.DeleteSingleAttachment(path)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": utils.Trans("cantDeleteFile", nil)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": utils.Trans("File delete Successfull", nil)})
}

func (ac *attachmentController) UploadSingleFile(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}
	err := context.Request.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		logger.LogInfo(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ownerID := context.GetInt64("user_id")
	requestParams := &commonModels.AccountRequstParams{
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
	}

	formdata := context.Request.MultipartForm
	fileHeaders := formdata.File["attachment"]
	attachmentType := context.Request.FormValue("attachment_type")
	check, typeOfAttachment := utils.ValidateAttachmentType(attachmentType)
	if !check {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidAttachmentType", nil)})
		return
	}
	if err != nil {
		logger.LogInfo(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidRequestParams", nil)})
		return
	}
	fileNames := strings.Split(context.Query("name"), ",")

	attachmentKey, err := utils.GenerateKeyHash()
	if err != nil {
		logger.LogInfo(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("ErrorWhileGeneratingHash", nil)})
		return
	}
	attachmentKey = attachmentKey + "_" + typeOfAttachment

	attachments, upErr := ac.attachmentService.UploadAttachments(requestParams, fileHeaders, fileNames, attachmentKey)
	if upErr != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": upErr})
		return
	}
	context.JSON(http.StatusOK, gin.H{"attachments": attachments})
	//return
}

func (ac *attachmentController) GetSingleAttachmentFile(context *gin.Context) {
	path := "attachments/" + context.Param("path")
	ac.attachmentService.GetSingleAttachment(context, path)
}

func (ac *attachmentController) FetchAttachments(context *gin.Context) {

	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	//accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		//accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	attachmentKey := context.Param("attachmentKey")
	attachments, err := ac.attachmentService.FetchAttachments(attachmentKey)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidRequestParams", nil)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"attachments": attachments})
}

func (ac *attachmentController) UploadMultipleAttachment(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	//accountId := context.GetUint64("account_id")
	var accountInfo commonModels.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModels.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		//accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	// err := context.Request.ParseMultipartForm(32 << 20) // 32MB
	// if err != nil {
	// 	logger.LogInfo(err)
	// 	context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// formData := context.Request.MultipartForm
	// fileHeaders := formData.File["attachment"]
	ownerID := context.GetInt64("user_id")
	accountId := context.GetUint64("account_id")
	requestParams := &commonModels.AccountRequstParams{
		CreatedBy: uint64(ownerID),
		AccountID: accountId,
	}
	reqBody := &models.UploadAttachmentRequestBody{}
	if err := context.BindJSON(reqBody); err != nil {
		logger.LogError(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("bindingError", nil)})
		return
	}
	attachmentType := reqBody.AttachmentType
	check, typeOfAttachment := utils.ValidateAttachmentType(attachmentType)
	if !check {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidAttachmentType", nil)})
		return
	}

	attachmentPaths := strings.Split(reqBody.Path, ",")
	attachmentNames := strings.Split(reqBody.Name, ",")

	if len(attachmentPaths) != len(attachmentNames) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("invalidRequestParams", nil)})
		return
	}

	attachmentKey, err := utils.GenerateKeyHash()
	if err != nil {
		logger.LogInfo(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("ErrorWhileGeneratingHash", nil)})
		return
	}
	attachmentKey = attachmentKey + "_" + typeOfAttachment

	attachments, upErr := ac.attachmentService.UploadMultipleAttachmentPath(requestParams, attachmentPaths, attachmentNames, attachmentKey)
	if upErr != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": upErr})
		return
	}

	context.JSON(http.StatusOK, gin.H{"attachments": attachments})
}
