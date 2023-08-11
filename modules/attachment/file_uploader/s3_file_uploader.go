package file_uploader

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"pi-inventory/common/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

type S3FileUploader struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

func (f *S3FileUploader) GetSession() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(f.Region),
			Credentials: credentials.NewStaticCredentials(
				f.AccessKeyID,
				f.SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (f *S3FileUploader) UploadSingleFile(directory, filename string, fileHeader *multipart.FileHeader) error {
	sess, err := f.GetSession()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	// fmt.Println("From s3")

	file_, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	body := bytes.NewReader(file_)

	defer file.Close()

	_, uploadErr := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(f.Bucket),
		Key:    aws.String(filename),
		Body:   body,
	})
	if uploadErr != nil {
		return uploadErr
	}

	return nil
}

func (f *S3FileUploader) GetSingleFile(context *gin.Context, filePath string, fileName string) {
	context.JSON(http.StatusOK, "not implemented yet")
}

func (f *S3FileUploader) DeleteSingleFile(filePath string) error {
	sess, err := f.GetSession()
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	out, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &f.Bucket,
		Key:    &filePath,
	})
	if err != nil {
		return err
	}
	logger.LogInfo(out.String())
	logger.LogInfo("deleting: ", filePath)

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: &f.Bucket,
		Key:    &filePath,
	})
	if err != nil {
		return err
	}

	return nil
}
