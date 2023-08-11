package service

import (
	"mime/multipart"
	"net/http"
	commonModels "pi-inventory/common/models"
	"pi-inventory/common/utils"
	"pi-inventory/errors"
	attachmentConst "pi-inventory/modules/attachment/consts"
	fileUploader "pi-inventory/modules/attachment/file_uploader"
	"pi-inventory/modules/attachment/models"
	"pi-inventory/modules/attachment/repository"
	"pi-inventory/modules/attachment/schema"
	"strings"

	"github.com/gin-gonic/gin"
)

type AttachmentServiceInterface interface {
	UploadAttachments(requestParams *commonModels.AccountRequstParams, files []*multipart.FileHeader, name []string, attachmentKey string) ([]*schema.Attachment, error)
	GetSingleAttachment(context *gin.Context, attachmentPath string)
	FetchAttachments(attachmentKey string) ([]*models.AttachmentCustomResponse, error)
	DeleteSingleAttachment(path string) error
	UploadMultipleAttachmentPath(requestParams *commonModels.AccountRequstParams, attachmentPaths []string, attachmentNames []string, attachmentKey string) ([]*schema.Attachment, error)
}

type attachmentService struct {
	attachmentRepository repository.AttachmentRepositoryInterface
	uploaderService      fileUploader.FileUploaderInterface
}

func NewAttachmentService(attachmentRepo repository.AttachmentRepositoryInterface, uploader fileUploader.FileUploaderInterface) *attachmentService {
	return &attachmentService{attachmentRepository: attachmentRepo, uploaderService: uploader}
}

func (as *attachmentService) UploadMultipleAttachmentPath(requestParams *commonModels.AccountRequstParams, attachmentPaths []string, attachmentNames []string, attachmentKey string) ([]*schema.Attachment, error) {
	// parentDirectory := "attachments"
	var attachments []*schema.Attachment

	// if len(attachmentNames) != len(fileHeaders) && len(fileHeaders) >= len(attachmentNames) {
	// 	unNamedFiles := len(fileHeaders) - len(attachmentNames)
	// 	for i := 0; i < unNamedFiles; i++ {
	// 		attachmentNames = append(attachmentNames, "")
	// 	}
	// }

	// if len(attachmentPaths) != len(fileHeaders) && len(fileHeaders) >= len(attachmentPaths) {
	// 	unPathFiles := len(fileHeaders) - len(attachmentPaths)
	// 	for i := 0; i < unPathFiles; i++ {
	// 		attachmentPaths = append(attachmentPaths, "")
	// 	}
	// }
	for i := 0; i < len(attachmentPaths); i++ {
		attachmentPath := attachmentPaths[i]
		attachmentName := attachmentNames[i]

		// err := as.uploaderService.UploadSingleFile(parentDirectory, attachmentPath, file)
		// if err != nil {
		// 	return nil, &errors.ApplicationError{
		// 		ErrorType:      errors.UnKnownErr,
		// 		TranslationKey: "unKnownError",
		// 		HttpCode:       http.StatusInternalServerError,
		// 	}
		// }

		// var attachmentName string
		// if index < len(attachmentNames) && strings.Compare(attachmentNames[index], "") != 0 {
		// 	attachmentName = attachmentNames[index]
		// } else {
		// 	nameWithExtension := strings.Split(attachmentPath, ".")
		// 	nameWithoutExtension := nameWithExtension[0]
		// 	attachmentName = nameWithoutExtension
		// }

		attachment, err := as.attachmentRepository.Store(requestParams.CreatedBy, attachmentPath, attachmentName, attachmentKey)
		if err != nil {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.UnKnownErr,
				TranslationKey: "unKnownError",
				HttpCode:       http.StatusInternalServerError,
			}
		}

		// attachment.Path = "/attachment/single/" + attachmentPath
		// if strings.Compare(attachment.Name, "") == 0 {
		// 	attachment.Name = attachmentPath
		// }

		attachments = append(attachments, attachment)

	}
	return attachments, nil
}

func (as *attachmentService) UploadAttachments(requestParams *commonModels.AccountRequstParams, files []*multipart.FileHeader, name []string, attachmentKey string) ([]*schema.Attachment, error) {
	parentDirectory := "attachments"
	var attachments []*schema.Attachment

	if len(name) != len(files) && len(files) >= len(name) {
		unNamedFiles := len(files) - len(name)
		for i := 0; i < unNamedFiles; i++ {
			name = append(name, "")
		}
	}

	for index, file := range files {
		attachmentPath := parentDirectory + "/" + utils.GetFileHashName(file.Filename)

		err := as.uploaderService.UploadSingleFile(parentDirectory, attachmentPath, file)
		if err != nil {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.UnKnownErr,
				TranslationKey: "unKnownError",
				HttpCode:       http.StatusInternalServerError,
			}
		}

		pathSegregation := strings.Split(attachmentPath, "/")
		fileName := pathSegregation[1]

		if strings.Compare(name[index], "") == 0 {
			nameWithExtension := strings.Split(fileName, ".")
			nameWithoutExtension := nameWithExtension[0]
			name[index] = nameWithoutExtension
		}

		attachment, err := as.attachmentRepository.Store(requestParams.CreatedBy, attachmentPath, name[index], attachmentKey)
		if err != nil {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.UnKnownErr,
				TranslationKey: "unKnownError",
				HttpCode:       http.StatusInternalServerError,
			}
		}
		//loop ends

		// Attachment Path modification
		attachment.Path = "/attachment/single/" + fileName
		if strings.Compare(attachment.Name, "") == 0 {
			attachment.Name = fileName
		}

		attachments = append(attachments, attachment)

	}

	return attachments, nil
}

func (as *attachmentService) GetSingleAttachment(context *gin.Context, attachmentPath string) {
	attachment, err := as.attachmentRepository.FindBy(attachmentConst.FieldPath, attachmentPath)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("fileNotFound", nil), "status": "error"})
		return
	}
	as.uploaderService.GetSingleFile(context, attachmentPath, attachment[0].Name)
}

func (as *attachmentService) DeleteSingleAttachment(attachmentPath string) error {
	return as.uploaderService.DeleteSingleFile(attachmentPath)
}

func (as *attachmentService) FetchAttachments(attachmentKey string) ([]*models.AttachmentCustomResponse, error) {
	attachments, err := as.attachmentRepository.FindBy(attachmentConst.FieldAttachmentKey, attachmentKey)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	resAttachments := make([]*models.AttachmentCustomResponse, 0)
	for _, attachment := range attachments {
		customAttachment := models.AttachmentCustomResponse{
			ID:   attachment.ID,
			Path: attachment.Path,
		}
		resAttachments = append(resAttachments, &customAttachment)
	}
	return resAttachments, nil
}
