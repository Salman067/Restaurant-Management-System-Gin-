package file_uploader

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"pi-inventory/common/utils"
)

type LocalFileUploader struct {
}

func (f LocalFileUploader) UploadSingleFile(directory, filename string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.MkdirAll(directory, 0700)
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}

func (f LocalFileUploader) GetSingleFile(context *gin.Context, filePath string, fileName string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": utils.Trans("fileNotFound", nil), "status": "error"})
		return
	}
	context.Header("X-File-Name", fileName)
	context.FileAttachment(filePath, fileName)
}

func (f LocalFileUploader) DeleteSingleFile(filePath string) error {
	return os.Remove(filePath)
}
