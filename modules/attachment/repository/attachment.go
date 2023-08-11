package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	attachmentConst "pi-inventory/modules/attachment/consts"
	"pi-inventory/modules/attachment/schema"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttachmentRepositoryInterface interface {
	Store(ownerID uint64, path string, name string, attachmentKey string) (*schema.Attachment, error)
	FindBy(field string, value any) ([]*schema.Attachment, error)
}

type attachmentRepository struct {
	Db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *attachmentRepository {
	return &attachmentRepository{Db: db}
}

func (ar *attachmentRepository) Store(ownerID uint64, path string, name string, attachmentKey string) (*schema.Attachment, error) {
	attachment := schema.Attachment{
		Name:          name,
		Path:          path,
		AttachmentKey: attachmentKey,
		OwnerID:       ownerID,
		CreatedAt:     time.Now(),
		CreatedBy:     1,
	}

	result := ar.Db.Clauses(clause.Returning{}).Create(&attachment)
	if result.Error != nil {
		logger.LogError(result.Error)
		return nil, result.Error
	}
	return &attachment, nil
}

func (ar *attachmentRepository) FindBy(field string, value any) ([]*schema.Attachment, error) {
	attachment := []*schema.Attachment{}
	query := fmt.Sprintf("%s = ?", field)
	logger.LogError(value)
	if err := ar.Db.Table(attachmentConst.AttachmentTable).Where(query, value).Find(&attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}
