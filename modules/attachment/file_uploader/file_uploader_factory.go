package file_uploader

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type FileUploaderInterface interface {
	UploadSingleFile(directory, filename string, fileHeader *multipart.FileHeader) error
	GetSingleFile(context *gin.Context, filePath string, fileName string)
	DeleteSingleFile(filePath string) error
}

func NewFileUploaderFactory() FileUploaderInterface {
	fileUploaderType := viper.GetString("FILE_UPLOADER")
	if fileUploaderType == "s3" {
		var s3Uploader = &S3FileUploader{
			AccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
			Region:          viper.GetString("AWS_REGION"),
			Bucket:          viper.GetString("BUCKET_NAME"),
		}
		return s3Uploader
	}

	return &LocalFileUploader{}
}
