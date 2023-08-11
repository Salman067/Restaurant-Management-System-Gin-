package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	lineItemConst "pi-inventory/modules/composite/consts"
	"pi-inventory/modules/composite/schema"

	"gorm.io/gorm"
)

type LineItemRepositoryInterface interface {
	Create(LineItems []*schema.LineItem) error
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.LineItem, error)
	Delete(LineItems []*schema.LineItem) error
}

type lineItemRepository struct {
	Db *gorm.DB
}

func NewLineItemRepository(db *gorm.DB) *lineItemRepository {
	return &lineItemRepository{Db: db}
}

// func (lr *lineItemRepository) FindAll(ownerID uint64, page *pagination.Page) ([]*schema.LineItem, int64, error) {
// 	var count int64
// 	lineItems := make([]*schema.LineItem, 0)
// 	if err := lr.Db.Table(lineItemConst.LineItemTable).Count(&count).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	baseQuery := lr.Db.Model(&schema.LineItem{}).Offset(page.Offset).Limit(page.Limit).Where("owner_id=?", ownerID).Order("id asc, id")
// 	if err := baseQuery.Find(&lineItems).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, count, err
// 		}
// 		return lineItems, count, nil
// 	}

// 	return lineItems, count, nil
// }

func (lr *lineItemRepository) Create(LineItem []*schema.LineItem) error {
	err := lr.Db.Table(lineItemConst.LineItemTable).Create(&LineItem).Error
	if err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
func (lr *lineItemRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.LineItem, error) {
	lineItem := make([]*schema.LineItem, 0)
	query := fmt.Sprintf("%s=? AND %s = ? ", lineItemConst.LineItemFieldAccountID, field)
	if err := lr.Db.Table(lineItemConst.LineItemTable).Where(query, requestParams.AccountID, value).Find(&lineItem).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return lineItem, nil
}
func (lr *lineItemRepository) Delete(LineItem []*schema.LineItem) error {
	err := lr.Db.Table(lineItemConst.LineItemTable).Save(&LineItem).Error
	if err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
