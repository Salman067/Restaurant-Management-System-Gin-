package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	compositeConst "pi-inventory/modules/composite/consts"
	"pi-inventory/modules/composite/models"
	"pi-inventory/modules/composite/schema"

	"gorm.io/gorm"
)

type CompositeRepositoryInterface interface {
	Create(composite *schema.Composite) (*schema.Composite, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.Composite, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.CompositeQueryParams) ([]*schema.Composite, int64, error)
	UpdateComposite(composite []*schema.Composite) (*schema.Composite, error)
	CompositeItemSummary(requestParams commonModel.AccountRequstParams) (int64, error)
}

type compositeRepository struct {
	Db *gorm.DB
}

func NewCompositeRepository(db *gorm.DB) *compositeRepository {
	return &compositeRepository{Db: db}
}

func (cr *compositeRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.CompositeQueryParams) ([]*schema.Composite, int64, error) {
	var count int64
	composites := make([]*schema.Composite, 0)
	query := fmt.Sprintf("%s=?", compositeConst.CompositeFieldAccountID)
	if err := cr.Db.Table(compositeConst.CompositeTable).Where(query, requestParams.AccountID).Count(&count).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	baseQuery := cr.Db.Model(&schema.Composite{}).
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit).
		Where(query, requestParams.AccountID).
		Order("created_at DESC")

	// Full text search
	//query := fmt.Sprintf("to_tsvector(%s || ' '|| %s || ' '|| %s || ' '|| %s || ' '|| %s) @@ to_tsquery(?)",
	//	compositeConst.CompositeFieldTitle, compositeConst.CompositeFieldTag, compositeConst.CompositeFieldDescription,
	//	compositeConst.CompositeFieldSellingPrice, compositeConst.CompositeFieldPurchasePrice)
	//if len(queryParamStruct.KeyWord) != 0 {
	//	if err := baseQuery.Where(query, queryParamStruct.KeyWord).Find(&composites).Error; err != nil {
	//		logger.LogError(err)
	//		return nil, 0, err
	//	}
	//}

	//Partial search
	searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
		compositeConst.CompositeFieldTitle, compositeConst.CompositeFieldTag, compositeConst.CompositeFieldDescription,
		compositeConst.CompositeFieldSellingPrice, compositeConst.CompositeFieldPurchasePrice)

	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		if err := baseQuery.Where(searchQ, keyword, keyword, keyword, keyword, keyword).Find(&composites).Count(&count).Error; err != nil {
			logger.LogError(err)
			return nil, 0, err
		}
	}

	if err := baseQuery.Find(&composites).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	return composites, count, nil
}

func (cr *compositeRepository) Create(Composite *schema.Composite) (*schema.Composite, error) {
	err := cr.Db.Table(compositeConst.CompositeTable).Create(&Composite).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return Composite, nil
}

func (cr *compositeRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.Composite, error) {
	composite := make([]*schema.Composite, 0)
	query := fmt.Sprintf("%s = ? AND %s=? ", compositeConst.CompositeFieldAccountID, field)
	if err := cr.Db.Table(compositeConst.CompositeTable).Where(query, requestParams.AccountID, value).First(&composite).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return composite, nil
}

func (cr *compositeRepository) UpdateComposite(composite []*schema.Composite) (*schema.Composite, error) {
	if err := cr.Db.Table(compositeConst.CompositeTable).Save(&composite).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return composite[0], nil
}

func (cr *compositeRepository) CompositeItemSummary(requestParams commonModel.AccountRequstParams) (int64, error) {
	var count int64
	query := fmt.Sprintf("%s=?", compositeConst.CompositeFieldAccountID)
	if err := cr.Db.Table(compositeConst.CompositeTable).Count(&count).Where(query, requestParams.AccountID).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return count, nil
}
