package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	"pi-inventory/common/models"
	commonModel "pi-inventory/common/models"
	stockConst "pi-inventory/modules/stock/consts"
	stockModel "pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/schema"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type StockRepositoryInterface interface {
	CreateStock(stock *schema.Stock) (*schema.Stock, error)
	FindAll(requestParams commonModel.AccountRequstParams, page *models.Page, queryParamStruct *stockModel.StockQueryParams) ([]*schema.Stock, int64, error)
	FindAllNew(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) ([]*stockModel.CustomQueryStockResponse, int64, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.Stock, error)
	UpdateStock(stock []*schema.Stock) (*schema.Stock, error)
	DeleteStock(requestParams commonModel.AccountRequstParams, ID uint64) (*schema.Stock, error)
	OutOfStock(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) ([]*stockModel.CustomQueryStockResponse, int64, error)
	LowOnStock(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) ([]*stockModel.CustomQueryStockResponse, int64, error)
	CountOfStockList(requestParams commonModel.AccountRequstParams, field string, value any, queryParamStruct *stockModel.StockQueryParams) (int64, error)
	CountOfItemList(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) (int64, error)
	CountOfExpiryDate(requestParams commonModel.AccountRequstParams) (int64, error)
	CreateStockActivityWhileStockCreatingOrUpdating(stockActivity *schema.StockActivity) (*schema.StockActivity, error)
}

type stockRepository struct {
	Db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	return &stockRepository{Db: db}
}

func (sr *stockRepository) CreateStock(stock *schema.Stock) (*schema.Stock, error) {
	if err := sr.Db.Clauses(clause.Returning{}).Table(stockConst.StockTable).Create(&stock).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stock, nil
}

func (sr *stockRepository) FindAll(requestParams commonModel.AccountRequstParams, page *models.Page, queryParamStruct *stockModel.StockQueryParams) ([]*schema.Stock, int64, error) {
	var count int64
	stocks := make([]*schema.Stock, 0)
	if err := sr.Db.Table(stockConst.StockTable).Count(&count).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	query := fmt.Sprintf("%s=?", stockConst.StockFieldAccountID)
	baseQuery := sr.Db.Model(&schema.Stock{}).
		Offset(page.Offset).
		Limit(page.Limit).
		Where(query, requestParams.AccountID).
		Order("created_at DESC")

	if queryParamStruct.Type == "" && queryParamStruct.Status == "" {
		baseQuery.Find(&stocks)
		return stocks, count, nil
	}
	if len(queryParamStruct.Type) != 0 {
		baseQuery = baseQuery.Where("stocks.type = ?", queryParamStruct.Type)
	}
	if len(queryParamStruct.Status) != 0 {
		baseQuery = baseQuery.Where("stocks.status = ?", queryParamStruct.Status)
	}
	if err := baseQuery.Find(&stocks).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.LogError(err)
			return nil, 0, err
		}
		return nil, 0, nil
	}

	return stocks, count, nil
}

func (sr *stockRepository) FindAllNew(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) ([]*stockModel.CustomQueryStockResponse, int64, error) {
	customStocks := make([]*stockModel.CustomQueryStockResponse, 0)
	count, err := sr.getCount(requestParams, queryParamStruct)
	if err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	baseQuery := sr.Db.Model(&schema.Stock{}).
		Select("stocks.*, units.title AS unit_title, categories.title AS category_title, warehouses.title AS location_title").
		Joins("LEFT JOIN units ON units.id = stocks.unit_id").
		Joins("LEFT JOIN categories ON categories.id = stocks.category_id").
		Joins("LEFT JOIN warehouses ON warehouses.id = stocks.location_id").
		Where("stocks.account_id=?", requestParams.AccountID).Where("is_deleted =?", false).
		Order("created_at DESC").
		Offset(requestParams.Page.Offset).Limit(requestParams.Page.Limit)

	if queryParamStruct.Type == "" && queryParamStruct.Status == "" && queryParamStruct.ListType == "" && queryParamStruct.KeyWord == "" {
		baseQuery.Find(&customStocks)
		return customStocks, count, nil
	}

	// Full text search
	//query := fmt.Sprintf("to_tsvector(%s || ' '|| %s || ' '|| %s || ' '|| %s || ' '|| stocks.%s) @@ to_tsquery(?)",
	//	stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.StockFieldDescription)
	//if len(queryParamStruct.KeyWord) != 0 {
	//	baseQuery = baseQuery.Where(query, queryParamStruct.KeyWord)
	//}

	// Partial search
	var searchQ string
	if queryParamStruct != nil && strings.Compare(queryParamStruct.Type, stockConst.StockInventoryType) == 0 {
		searchQ = fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
			stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.WarehouseFieldTitle)
	} else {
		searchQ = fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
			stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.UnitTableFieldTitle)
	}

	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		baseQuery = baseQuery.Where(searchQ, keyword, keyword, keyword, keyword, keyword).Count(&count)
	}

	if len(queryParamStruct.Type) != 0 {
		baseQuery = baseQuery.Where("stocks.type = ?", queryParamStruct.Type)
	}
	if len(queryParamStruct.ListType) != 0 {
		if queryParamStruct.ListType == "expired" {
			today := time.Now().Format("2006-01-02") // Get the current date
			baseQuery = baseQuery.Where("DATE(stocks.expiry_date) < ?", today)
		}
	}
	if len(queryParamStruct.Status) != 0 {
		baseQuery = baseQuery.Where("stocks.status = ?", queryParamStruct.Status)
	}
	if err := baseQuery.Find(&customStocks).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	return customStocks, count, nil
}

func (sr *stockRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.Stock, error) {
	stock := make([]*schema.Stock, 0)
	query := fmt.Sprintf("%s = ? ", field)
	if err := sr.Db.Table(stockConst.StockTable).
		Where("account_id = ?", requestParams.AccountID).
		Where("is_deleted =?", false).
		Where(query, value).Find(&stock).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stock, nil
}

func (sr *stockRepository) UpdateStock(stock []*schema.Stock) (*schema.Stock, error) {
	if err := sr.Db.Table(stockConst.StockTable).Save(&stock).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stock[0], nil
}

func (sr *stockRepository) DeleteStock(requestParams commonModel.AccountRequstParams, ID uint64) (*schema.Stock, error) {
	stock := &schema.Stock{}
	query := fmt.Sprintf("%s=? AND %s=?", stockConst.StockFieldAccountID, stockConst.StockFieldID)
	err := sr.Db.Table(stockConst.StockTable).Where(query, requestParams.AccountID, ID).First(&stock).Update("is_deleted", true).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stock, nil
}

func (sr *stockRepository) OutOfStock(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) ([]*stockModel.CustomQueryStockResponse, int64, error) {
	var count int64
	var stockList []*stockModel.CustomQueryStockResponse

	baseQuery := sr.Db.Model(&schema.Stock{}).
		Select("stocks.*, units.title AS unit_title, categories.title AS category_title, warehouses.title AS location_title").
		Joins("LEFT JOIN units ON units.id = stocks.unit_id").
		Joins("LEFT JOIN categories ON categories.id = stocks.category_id").
		Joins("LEFT JOIN warehouses ON warehouses.id = stocks.location_id").
		Where("stocks.account_id=?", requestParams.AccountID).
		Where("is_deleted =?", false).
		Order("id DESC, id")

	searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
		stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.WarehouseFieldTitle)

	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		baseQuery = baseQuery.Where(searchQ, keyword, keyword, keyword, keyword, keyword).Count(&count)
	}

	query := fmt.Sprintf("stocks.stock_qty = 0 AND %s = ?", stockConst.StockFieldType)
	if err := baseQuery.Table(stockConst.StockTable).
		Where(query, stockConst.TypeInventory).
		Count(&count).
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit).
		Find(&stockList).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	return stockList, count, nil
}

func (sr *stockRepository) LowOnStock(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) ([]*stockModel.CustomQueryStockResponse, int64, error) {
	var count int64
	var stockList []*stockModel.CustomQueryStockResponse

	baseQuery := sr.Db.Model(&schema.Stock{}).
		Select("stocks.*, units.title AS unit_title, categories.title AS category_title, warehouses.title AS location_title").
		Joins("LEFT JOIN units ON units.id = stocks.unit_id").
		Joins("LEFT JOIN categories ON categories.id = stocks.category_id").
		Joins("LEFT JOIN warehouses ON warehouses.id = stocks.location_id").
		Where("stocks.account_id=?", requestParams.AccountID).
		Where("is_deleted =?", false).
		Order("id DESC, id")

	searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
		stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.WarehouseFieldTitle)

	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		baseQuery = baseQuery.Where(searchQ, keyword, keyword, keyword, keyword, keyword).Count(&count)
	}

	query := fmt.Sprintf("stocks.reorder_qty > stock_qty AND stock_qty != 0 AND %s = ?", stockConst.StockFieldType)
	if err := baseQuery.Table(stockConst.StockTable).
		Where(query, stockConst.TypeInventory).
		Count(&count).
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit).
		Find(&stockList).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	return stockList, count, nil
}

func (sr *stockRepository) getCount(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) (int64, error) {
	var count int64
	baseQuery := sr.Db.Model(&schema.Stock{}).
		Where("stocks.account_id=?", requestParams.AccountID).
		Where("is_deleted =?", false)

	query := fmt.Sprintf("to_tsvector(%s || ' '|| %s || ' '|| %s || ' '|| %s || ' '|| %s) @@ to_tsquery(?)",
		stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.StockFieldDescription)
	if len(queryParamStruct.KeyWord) != 0 {
		baseQuery = baseQuery.Where(query, queryParamStruct.KeyWord)
	}
	if len(queryParamStruct.Type) != 0 {
		baseQuery = baseQuery.Where("stocks.type = ?", queryParamStruct.Type)
	}
	if len(queryParamStruct.ListType) != 0 {
		if queryParamStruct.ListType == "expired" {
			today := time.Now().Format("2006-01-02") // Get the current date
			baseQuery = baseQuery.Where("DATE(stocks.expiry_date) < ?", today)
		}
	}
	if len(queryParamStruct.Status) != 0 {
		baseQuery = baseQuery.Where("stocks.status = ?", queryParamStruct.Status)
	}
	if queryParamStruct.Type == stockConst.TypeInventory {
		if err := sr.Db.Table(stockConst.StockTable).Where("stocks.type = ? AND account_id=?", queryParamStruct.Type, requestParams.AccountID).Count(&count).Error; err != nil {
			logger.LogError(err)
			return 0, err
		}
	}
	if err := baseQuery.Count(&count).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return count, nil
}

func (sr *stockRepository) CountOfStockList(requestParams commonModel.AccountRequstParams, field string, value any, queryParamStruct *stockModel.StockQueryParams) (int64, error) {
	var count int64

	query := fmt.Sprintf("%s = ? AND %s = ?", stockConst.StockFieldAccountID, field)
	baseQuery := sr.Db.Table(stockConst.StockTable).
		Where(query, requestParams.AccountID, value).Where("is_deleted =?", false)

	searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR stocks.%s ILIKE ?",
		stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.StockFieldDescription)

	if queryParamStruct != nil && len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		baseQuery = baseQuery.Where(searchQ, keyword, keyword, keyword, keyword, keyword).Count(&count)
	}

	if err := baseQuery.Count(&count).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return count, nil
}

func (sr *stockRepository) CountOfItemList(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.StockQueryParams) (int64, error) {
	var count int64
	query := fmt.Sprintf("%s = ?", stockConst.StockFieldAccountID)
	baseQuery := sr.Db.Table(stockConst.StockTable).Where(query, requestParams.AccountID).Where("is_deleted =?", false)

	searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR stocks.%s ILIKE ?",
		stockConst.StockFieldName, stockConst.StockFieldSKU, stockConst.StockFieldSellingPrice, stockConst.StockFieldPurchasePrice, stockConst.StockFieldDescription)

	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		baseQuery = baseQuery.Where(searchQ, keyword, keyword, keyword, keyword, keyword).Count(&count)
	}

	if err := baseQuery.Count(&count).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return count, nil
}

func (sr *stockRepository) CountOfExpiryDate(requestParams commonModel.AccountRequstParams) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02") // Get the current date
	query := fmt.Sprintf("%s = ? AND DATE(stocks.expiry_date) < ?", stockConst.StockFieldAccountID)
	if err := sr.Db.Table(stockConst.StockTable).
		Where(query, requestParams.AccountID, today).Where("is_deleted =?", false).Count(&count).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return count, nil
}

func (sr *stockRepository) CreateStockActivityWhileStockCreatingOrUpdating(stockActivity *schema.StockActivity) (*schema.StockActivity, error) {
	if err := sr.Db.Clauses(clause.Returning{}).Table(stockConst.StockActivityTable).Create(&stockActivity).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return stockActivity, nil
}
