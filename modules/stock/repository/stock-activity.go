package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	stockActivityConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"

	"pi-inventory/modules/stock/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockActivityRepositoryInterface interface {
	CreateStockActivity(stockActivity *schema.StockActivity) (*schema.StockActivity, error)
	FindAll(requestParams commonModel.AccountRequstParams, paramStruct *models.StockActivityQueryParams) ([]*models.CustomStockActivityResponse, int64, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.StockActivity, error)
}

type stockActivityRepository struct {
	Db *gorm.DB
}

func NewStockActivityRepository(db *gorm.DB) *stockActivityRepository {
	return &stockActivityRepository{Db: db}
}

func (sar *stockActivityRepository) CreateStockActivity(stockActivity *schema.StockActivity) (*schema.StockActivity, error) {
	if err := sar.Db.Clauses(clause.Returning{}).Table(stockActivityConst.StockActivityTable).Create(&stockActivity).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stockActivity, nil
}

func (sar *stockActivityRepository) FindAll(requestParams commonModel.AccountRequstParams, paramStruct *models.StockActivityQueryParams) ([]*models.CustomStockActivityResponse, int64, error) {
	var count int64
	customStockActivities := make([]*models.CustomStockActivityResponse, 0)
	baseQuery := sar.Db.Table("stock_activities").
		Select("stock_activities.*, purposes.title AS purpose_title, stocks.name AS stock_title").
		Joins("LEFT JOIN purposes ON purposes.id = stock_activities.purpose_id").
		Joins("LEFT JOIN stocks ON stocks.id = stock_activities.stock_id").
		Order("id DESC, id").
		Where("stock_activities.account_id = ?", requestParams.AccountID).Where("stock_activities.mode != ?", "create")

	if paramStruct.StockID == 0 && paramStruct.Keyword == "" {
		baseQuery.Count(&count).Offset(requestParams.Page.Offset).Limit(requestParams.Page.Limit).Find(&customStockActivities)
		return customStockActivities, count, nil
	} else {
		// Partial search
		if len(paramStruct.Keyword) > 0 {
			keyword := fmt.Sprintf("%%%s%%", paramStruct.Keyword) // Add wildcards to the keyword
			searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ?",
				stockActivityConst.StockActivityFieldPurposeTitle, stockActivityConst.StockActivityFieldStockTitle)
			baseQuery = baseQuery.Where(searchQ, keyword, keyword).Offset(requestParams.Page.Offset).Limit(requestParams.Page.Limit)
		}
		if paramStruct.StockID != 0 {
			baseQuery = baseQuery.Where("stock_activities.stock_id = ?", paramStruct.StockID)
		}

		if err := baseQuery.Count(&count).Offset(requestParams.Page.Offset).Limit(requestParams.Page.Limit).Find(&customStockActivities).Error; err != nil {
			logger.LogError(err)
			return nil, 0, err
		}
		return customStockActivities, count, nil
	}

}

func (sar *stockActivityRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.StockActivity, error) {
	stockActivity := make([]*schema.StockActivity, 0)
	query := fmt.Sprintf("%s = ? ", field)
	if err := sar.Db.Table(stockActivityConst.StockActivityTable).Where("account_id = ?", requestParams.AccountID).Where(query, value).First(&stockActivity).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stockActivity, nil
}
