package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	attachmentModel "pi-inventory/modules/attachment/models"
	attachmentModuleService "pi-inventory/modules/attachment/service"
	"pi-inventory/modules/stock/cache"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
	"pi-inventory/modules/stock/schema"
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type StockServiceInterface interface {
	CreateStock(requestParams commonModel.AccountRequstParams, reqStock *models.AddStockRequestBody, context context.Context) (*schema.Stock, error)
	FindAll(requestParams commonModel.AccountRequstParams, page *commonModel.Page, queryParamStruct *models.StockQueryParams) ([]*models.SingleStockResponse, *commonModel.PageResponse, error)
	FindAllNew(requestParams commonModel.AccountRequstParams, queryParamStruct *models.StockQueryParams) ([]*models.SingleStockResponse, *commonModel.PageResponse, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.SingleStockResponseView, error)
	UpdateStock(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateStockRequestBody) (*schema.Stock, error)
	DeleteStock(requestParams commonModel.AccountRequstParams, ID uint64) (*schema.Stock, error)
	StockSummary(requestParams commonModel.AccountRequstParams, queryParamStruct *models.StockQueryParams) (*models.StockSummaryResponse, error)
}
type stockService struct {
	stockRepository      repository.StockRepositoryInterface
	attachmentService    attachmentModuleService.AttachmentServiceInterface
	categoryService      CategoryServiceInterface
	unitService          UnitServiceInterface
	stockCacheRepository cache.StockCacheRepositoryInterface
}

func NewStockService(stockRepo repository.StockRepositoryInterface,
	attachmentSrv attachmentModuleService.AttachmentServiceInterface,
	categoryService CategoryServiceInterface,
	unitService UnitServiceInterface,
	stockCacheRepository cache.StockCacheRepositoryInterface) *stockService {
	return &stockService{
		stockRepository:      stockRepo,
		attachmentService:    attachmentSrv,
		categoryService:      categoryService,
		unitService:          unitService,
		stockCacheRepository: stockCacheRepository,
	}
}

func (ss *stockService) CreateStock(requestParams commonModel.AccountRequstParams, reqStock *models.AddStockRequestBody, context context.Context) (*schema.Stock, error) {

	err := reqStock.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	stocks, err := ss.stockRepository.FindBy(requestParams, stockConst.StockFieldName, reqStock.Name)

	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "UnknownError",
			HttpCode:       http.StatusInternalServerError}
	}
	if len(stocks) > 0 {
		logger.LogError("A stock with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "item",
			},
			HttpCode: http.StatusBadRequest,
		}
	}

	newStock := &schema.Stock{}
	setRequestStockFieldValues(newStock, reqStock, requestParams)

	if reqStock.CategoryID != 0 {
		_, err = ss.categoryService.FindByID(requestParams, reqStock.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	if reqStock.UnitID != 0 {
		_, err = ss.unitService.FindByID(requestParams, reqStock.UnitID)
		if err != nil {
			return nil, err
		}
	}

	if reqStock.TrackInventory && reqStock.Status == stockConst.StockInactiveStatus {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "inactiveItem",
			TranslationParams: map[string]interface{}{
				"field": "inactive items",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	stock, err := ss.stockRepository.CreateStock(newStock)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Stock",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	//Storing products in cache if that product is maintained as stock
	if stock.TrackInventory {
		key := "product_" + requestParams.AccountSlug
		byteStock, err := json.Marshal(stock)
		if err != nil {
			logger.LogError("Stock is not marshal while storing in cache")
		}

		stringStockUUID := fmt.Sprint(stock.StockUUID)
		value := map[string]interface{}{
			stringStockUUID: string(byteStock),
		}

		err = ss.stockCacheRepository.Set(context, key, value)
		if err != nil {
			logger.LogError("Stock is not set in the cache")
		}

	}

	if reqStock.CategoryID != 0 {
		updateCategoryReq := &models.UpdateCategoryRequestBody{
			IsUsed: true,
		}
		_, err = ss.categoryService.UpdateCategory(requestParams, stock.CategoryID, updateCategoryReq)
		if err != nil {
			return nil, err
		}
	}

	if reqStock.UnitID != 0 {
		updateUnitReq := &models.UpdateUnitRequestBody{
			IsUsed: true,
		}
		_, err = ss.unitService.UpdateUnit(requestParams, stock.UnitID, updateUnitReq)
		if err != nil {
			return nil, err
		}

	}

	if reqStock.TrackInventory {
		stockActivityRequest := getStockActivityRequest(stock, requestParams)

		_, err = ss.stockRepository.CreateStockActivityWhileStockCreatingOrUpdating(stockActivityRequest)
		if err != nil {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.Unsuccessfull,
				TranslationKey: "createUnsuccessfull",
				TranslationParams: map[string]interface{}{
					"field": "Stock activity",
				},
				HttpCode: http.StatusInternalServerError,
			}
		}
	}

	return stock, nil
}

func (ss *stockService) FindAll(requestParams commonModel.AccountRequstParams, page *commonModel.Page, queryParamStruct *models.StockQueryParams) ([]*models.SingleStockResponse, *commonModel.PageResponse, error) {
	stocks, recordCount, err := ss.stockRepository.FindAll(requestParams, page, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	responseStocks := setStockResponse(stocks)
	for i, singleStock := range stocks {
		attachments, err := ss.attachmentService.FetchAttachments(singleStock.AttachmentKey)
		if err != nil {
			return nil, nil, err
		}
		setAttachments(responseStocks[i], attachments)
	}
	pageInfo := commonModel.PageResponse{
		Offset: page.Offset,
		Limit:  page.Limit,
		Count:  int(recordCount),
	}

	return responseStocks, &pageInfo, nil
}

func (ss *stockService) FindAllNew(requestParams commonModel.AccountRequstParams, queryParamStruct *models.StockQueryParams) ([]*models.SingleStockResponse, *commonModel.PageResponse, error) {
	stocks, recordCount, err := ss.stockRepository.FindAllNew(requestParams, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	responseStocks := setCustomStockResponse(stocks)

	for i, singleStock := range stocks {
		attachments, err := ss.attachmentService.FetchAttachments(singleStock.AttachmentKey)
		if err != nil {
			return nil, nil, err
		}
		setAttachments(responseStocks[i], attachments)
	}
	pageInfo := commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}

	return responseStocks, &pageInfo, nil
}

func (ss *stockService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.SingleStockResponseView, error) {
	stocks, err := ss.stockRepository.FindBy(requestParams, stockConst.StockFieldID, ID)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "UnknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	if len(stocks) == 0 {
		logger.LogError("StockIDNotFound")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.RecordNotFound,
			TranslationKey: "idNotFound",
			TranslationParams: map[string]interface{}{
				"field": "stock",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	responseStocks := setStockResponseView(stocks)
	attachments, err := ss.attachmentService.FetchAttachments(stocks[0].AttachmentKey)
	if err != nil {
		return nil, err
	}
	setAttachmentsView(responseStocks[0], attachments)
	return responseStocks[0], nil
}

func (ss *stockService) UpdateStock(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateStockRequestBody) (*schema.Stock, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	stock, err := ss.stockRepository.FindBy(requestParams, stockConst.StockFieldID, ID)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "UnknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	if len(stock) == 0 {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.RecordNotFound,
			TranslationKey: "idNotFound",
			TranslationParams: map[string]interface{}{
				"field": "stock",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	if reqBody.CategoryID != 0 {
		_, err = ss.categoryService.FindByID(requestParams, reqBody.CategoryID)
		if err != nil {
			return nil, err
		}
	}
	if reqBody.UnitID != 0 {
		_, err = ss.unitService.FindByID(requestParams, reqBody.UnitID)
		if err != nil {
			return nil, err
		}
	}

	previousTrackInventoryValue := stock[0].TrackInventory

	checkedStock := updateStockField(reqBody, stock)
	if checkedStock[0].TrackInventory && checkedStock[0].Status == stockConst.StockInactiveStatus {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "inactiveItem",
			TranslationParams: map[string]interface{}{
				"field": "inactive items",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	stocks, err := ss.stockRepository.UpdateStock(checkedStock)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Stock",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	//Storing products in cache if that product is maintained as stock
	if stocks.TrackInventory {
		key := "product_" + requestParams.AccountSlug
		byteStock, err := json.Marshal(stocks)
		if err != nil {
			logger.LogError("Stock is not marshal while storing in cache")
		}

		stringStockID := fmt.Sprint(stocks.ID)
		value := map[string]interface{}{
			stringStockID: string(byteStock),
		}

		err = ss.stockCacheRepository.Set(context.Background(), key, value)
		if err != nil {
			logger.LogError("Stock is not set in the cache")
		}
	}

	if reqBody.CategoryID != 0 {
		updateCategoryReq := &models.UpdateCategoryRequestBody{
			IsUsed: true,
		}
		_, err = ss.categoryService.UpdateCategory(requestParams, stocks.CategoryID, updateCategoryReq)
		if err != nil {
			return nil, err
		}
	}

	if reqBody.UnitID != 0 {
		updateUnitReq := &models.UpdateUnitRequestBody{
			IsUsed: true,
		}
		_, err = ss.unitService.UpdateUnit(requestParams, stocks.UnitID, updateUnitReq)
		if err != nil {
			return nil, err
		}
	}

	// Checking for the stock to stock conversion or, it's a normal stock update
	// In case of conversion, it should track stock creation history in stock Activity schema
	// In case of normal update operation, it will not be considered as stock tracking history
	if previousTrackInventoryValue != stocks.TrackInventory {
		stockActivityRequest := getStockActivityRequest(stocks, requestParams)

		_, err = ss.stockRepository.CreateStockActivityWhileStockCreatingOrUpdating(stockActivityRequest)
		if err != nil {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.Unsuccessfull,
				TranslationKey: "createUnsuccessfull",
				TranslationParams: map[string]interface{}{
					"field": "Stock activity",
				},
				HttpCode: http.StatusInternalServerError,
			}
		}
	}

	return stocks, nil
}

func (ss *stockService) DeleteStock(requestParams commonModel.AccountRequstParams, ID uint64) (*schema.Stock, error) {

	stocks, err := ss.stockRepository.FindBy(requestParams, stockConst.StockFieldID, ID)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	if len(stocks) == 0 {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.RecordNotFound,
			TranslationKey: "idNotFound",
			TranslationParams: map[string]interface{}{
				"field": "stock",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	//Check if it is a stock or not before Delete
	if stocks[0].RecordType == "stock" {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteForbidden",
			TranslationParams: map[string]interface{}{
				"field": "Stock",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	stock, err := ss.stockRepository.DeleteStock(requestParams, ID)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "deleteUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Stock",
			},
			HttpCode: http.StatusInternalServerError}
	}
	return stock, nil

}

func (ss *stockService) StockSummary(requestParams commonModel.AccountRequstParams, queryParamStruct *models.StockQueryParams) (*models.StockSummaryResponse, error) {
	outOfStockList, outOfStockCount, err := ss.stockRepository.OutOfStock(requestParams, queryParamStruct)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	outOfStockPageInfo := &commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(outOfStockCount),
	}
	var responseOutOfStockList []*models.SingleStockResponse
	for _, singleStock := range outOfStockList {
		attachments, err := ss.attachmentService.FetchAttachments(singleStock.AttachmentKey)
		if err != nil {
			return nil, err
		}
		singleStockResponse := setStockSummeryListFields(singleStock, attachments)
		responseOutOfStockList = append(responseOutOfStockList, singleStockResponse)
	}

	lowOnStockList, lowOnStockCount, err := ss.stockRepository.LowOnStock(requestParams, queryParamStruct)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	lowOnStockPageInfo := &commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(lowOnStockCount),
	}
	var responseLowOnStockList []*models.SingleStockResponse
	for _, singleStock := range lowOnStockList {
		attachments, err := ss.attachmentService.FetchAttachments(singleStock.AttachmentKey)
		if err != nil {
			return nil, err
		}
		singleStockResponse := setStockSummeryListFields(singleStock, attachments)
		responseLowOnStockList = append(responseLowOnStockList, singleStockResponse)
	}

	countOfStockList, err := ss.stockRepository.CountOfStockList(requestParams, stockConst.StockFieldType, stockConst.TypeInventory, queryParamStruct)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}

	countOfItemList, err := ss.stockRepository.CountOfItemList(requestParams, queryParamStruct)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}

	countOfExpiryDate, err := ss.stockRepository.CountOfExpiryDate(requestParams)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	responseSummaryStock := &models.StockSummaryResponse{
		CountOfStockList:   countOfStockList,
		CountOfItemList:    countOfItemList,
		CountOfExpiryDate:  countOfExpiryDate,
		OutOfStock:         outOfStockCount,
		LowOnStock:         lowOnStockCount,
		OutOfStockPageInfo: outOfStockPageInfo,
		OutOfStocks:        responseOutOfStockList,
		LowOnStockPageInfo: lowOnStockPageInfo,
		LowOnStocks:        responseLowOnStockList,
	}

	if responseOutOfStockList == nil {
		responseSummaryStock.OutOfStocks = make([]*models.SingleStockResponse, 0)
	}
	if responseLowOnStockList == nil {
		responseSummaryStock.LowOnStocks = make([]*models.SingleStockResponse, 0)
	}

	return responseSummaryStock, nil
}

func setResponseStockFieldValues(instance *schema.Stock, resStock *models.SingleStockResponse) {

	resStock.ID = instance.ID
	resStock.OwnerID = instance.OwnerID
	resStock.AccountID = instance.AccountID
	resStock.Name = instance.Name
	resStock.Type = instance.Type
	resStock.SKU = instance.SKU
	resStock.Status = instance.Status
	resStock.Description = instance.Description
	resStock.SellingPrice = commonModel.Money{
		Amount:   instance.SellingPrice,
		Currency: instance.SellingPriceCurrency,
	}.ConvertFloatMoney()
	resStock.PurchasePrice = commonModel.Money{
		Amount:   instance.PurchasePrice,
		Currency: instance.PurchasePriceCurrency,
	}.ConvertFloatMoney()
	resStock.StockQty = instance.StockQty
	resStock.ReorderQty = instance.ReorderQty
	resStock.RecordType = instance.RecordType
	resStock.AsOfDate = instance.AsOfDate
	resStock.PurchaseDate = instance.PurchaseDate
	resStock.ExpiryDate = instance.ExpiryDate
	resStock.SupplierID = instance.SupplierID
}

func setStockResponse(list []*schema.Stock) []*models.SingleStockResponse {

	stocks := make([]*models.SingleStockResponse, 0)
	for _, singleInstance := range list {

		stock := models.SingleStockResponse{}
		setResponseStockFieldValues(singleInstance, &stock)
		stocks = append(stocks, &stock)

	}

	return stocks
}

func setCustomResponseStockFieldValues(instance *models.CustomQueryStockResponse, resStock *models.SingleStockResponse) {

	resStock.ID = instance.ID
	resStock.OwnerID = instance.OwnerID
	resStock.AccountID = instance.AccountID
	resStock.Name = instance.Name
	resStock.Type = instance.Type
	resStock.SKU = instance.SKU
	resStock.Status = instance.Status
	resStock.Description = instance.Description
	resStock.SellingPrice = commonModel.Money{
		Amount:   instance.SellingPrice,
		Currency: instance.SellingPriceCurrency,
	}.ConvertFloatMoney()
	resStock.PurchasePrice = commonModel.Money{
		Amount:   instance.PurchasePrice,
		Currency: instance.PurchasePriceCurrency,
	}.ConvertFloatMoney()
	resStock.StockQty = instance.StockQty
	resStock.ReorderQty = instance.ReorderQty
	resStock.RecordType = instance.RecordType
	resStock.AsOfDate = instance.AsOfDate
	resStock.PurchaseDate = instance.PurchaseDate
	resStock.ExpiryDate = instance.ExpiryDate
	resStock.UnitTitle = instance.UnitTitle
	resStock.CategoryTitle = instance.CategoryTitle
	resStock.LocationTitle = instance.LocationTitle
	resStock.SupplierID = instance.SupplierID
	resStock.SaleTaxID = instance.SaleTaxID
	resStock.PurchaseTaxID = instance.PurchaseTaxID
}

func setCustomStockResponse(list []*models.CustomQueryStockResponse) []*models.SingleStockResponse {

	stocks := make([]*models.SingleStockResponse, 0)
	for _, singleInstance := range list {

		stock := models.SingleStockResponse{}
		setCustomResponseStockFieldValues(singleInstance, &stock)
		stocks = append(stocks, &stock)

	}

	return stocks
}
func setStockSummeryListFields(instance *models.CustomQueryStockResponse, attachments []*attachmentModel.AttachmentCustomResponse) *models.SingleStockResponse {
	resStock := &models.SingleStockResponse{
		ID:          instance.ID,
		OwnerID:     instance.OwnerID,
		AccountID:   instance.AccountID,
		Name:        instance.Name,
		Type:        instance.Type,
		SKU:         instance.SKU,
		Status:      instance.Status,
		Description: instance.Description,
		SellingPrice: commonModel.Money{
			Amount:   instance.SellingPrice,
			Currency: instance.SellingPriceCurrency,
		}.ConvertFloatMoney(),
		PurchasePrice: commonModel.Money{
			Amount:   instance.PurchasePrice,
			Currency: instance.PurchasePriceCurrency,
		}.ConvertFloatMoney(),
		StockQty:      instance.StockQty,
		ReorderQty:    instance.ReorderQty,
		RecordType:    instance.RecordType,
		AsOfDate:      instance.AsOfDate,
		PurchaseDate:  instance.PurchaseDate,
		ExpiryDate:    instance.ExpiryDate,
		UnitTitle:     instance.UnitTitle,
		CategoryTitle: instance.CategoryTitle,
		LocationTitle: instance.LocationTitle,
		SupplierID:    instance.SupplierID,
		SaleTaxID:     instance.SaleTaxID,
		PurchaseTaxID: instance.PurchaseTaxID,
		Attachments:   attachments,
	}
	return resStock
}

func setStockResponseView(list []*schema.Stock) []*models.SingleStockResponseView {

	stocks := make([]*models.SingleStockResponseView, 0)
	for _, singleInstance := range list {

		stock := models.SingleStockResponseView{}
		setResponseStockFieldValuesView(singleInstance, &stock)
		stocks = append(stocks, &stock)

	}

	return stocks
}

func setResponseStockFieldValuesView(instance *schema.Stock, resStock *models.SingleStockResponseView) {

	resStock.ID = instance.ID
	resStock.OwnerID = instance.OwnerID
	resStock.AccountID = instance.AccountID
	resStock.Name = instance.Name
	resStock.Type = instance.Type
	resStock.SKU = instance.SKU
	resStock.Status = instance.Status
	resStock.Description = instance.Description
	resStock.SellingPrice = commonModel.Money{
		Amount:   instance.SellingPrice,
		Currency: instance.SellingPriceCurrency,
	}.ConvertFloatMoney()
	resStock.PurchasePrice = commonModel.Money{
		Amount:   instance.PurchasePrice,
		Currency: instance.PurchasePriceCurrency,
	}.ConvertFloatMoney()
	resStock.StockQty = instance.StockQty
	resStock.ReorderQty = instance.ReorderQty
	resStock.RecordType = instance.RecordType
	resStock.AsOfDate = instance.AsOfDate
	resStock.PurchaseDate = instance.PurchaseDate
	resStock.ExpiryDate = instance.ExpiryDate
	resStock.UnitID = instance.UnitID
	resStock.CategoryID = instance.CategoryID
	resStock.LocationID = instance.LocationID
	resStock.SupplierID = instance.SupplierID
	resStock.SaleTaxID = instance.SaleTaxID
	resStock.PurchaseTaxID = instance.PurchaseTaxID
}

func setAttachments(responseStock *models.SingleStockResponse, attachments []*attachmentModel.AttachmentCustomResponse) {
	responseStock.Attachments = append(responseStock.Attachments, attachments...)
}

func setAttachmentsView(responseStock *models.SingleStockResponseView, attachments []*attachmentModel.AttachmentCustomResponse) {
	responseStock.Attachments = append(responseStock.Attachments, attachments...)
}

func updateStockField(reqBody *models.UpdateStockRequestBody, stock []*schema.Stock) []*schema.Stock {
	if reqBody.Name != "" {
		stock[0].Name = reqBody.Name
	}
	if reqBody.SKU != "" {
		stock[0].SKU = reqBody.SKU
	}
	if reqBody.Status != "" {
		stock[0].Status = reqBody.Status
	}
	if reqBody.Description != "" {
		stock[0].Description = reqBody.Description
	}
	if reqBody.SellingPrice.Amount.GreaterThanOrEqual(decimal.Zero) {
		stock[0].SellingPrice = reqBody.SellingPrice.Amount
	}
	if reqBody.PurchasePrice.Amount.GreaterThanOrEqual(decimal.Zero) {
		stock[0].PurchasePrice = reqBody.PurchasePrice.Amount
	}
	if reqBody.SellingPrice.Currency != "" {
		stock[0].SellingPriceCurrency = reqBody.SellingPrice.Currency
	}
	if reqBody.PurchasePrice.Currency != "" {
		stock[0].PurchasePriceCurrency = reqBody.PurchasePrice.Currency
	}
	if reqBody.AttachmentKey != "" {
		stock[0].AttachmentKey = reqBody.AttachmentKey
	}
	if reqBody.UnitID != 0 {
		stock[0].UnitID = reqBody.UnitID
	}
	if reqBody.CategoryID != 0 {
		stock[0].CategoryID = reqBody.CategoryID
	}
	if reqBody.LocationID != 0 {
		stock[0].LocationID = reqBody.LocationID
	}
	if len(reqBody.SupplierID) != 0 {
		stock[0].SupplierID = reqBody.SupplierID
	}
	if reqBody.AsOfDate != nil {
		stock[0].AsOfDate = reqBody.AsOfDate
	}
	if reqBody.PurchaseDate != nil {
		stock[0].PurchaseDate = reqBody.PurchaseDate
	}
	if reqBody.ExpiryDate != nil {
		stock[0].ExpiryDate = reqBody.ExpiryDate
	}
	if len(reqBody.SaleTaxID) != 0 {
		stock[0].SaleTaxID = reqBody.SaleTaxID
	}
	if len(reqBody.PurchaseTaxID) != 0 {
		stock[0].PurchaseTaxID = reqBody.PurchaseTaxID
	}
	if reqBody.TrackInventory {
		stock[0].StockQty = reqBody.StockQty
		stock[0].ReorderQty = reqBody.ReorderQty
	}
	stock[0].TrackInventory = reqBody.TrackInventory
	stock[0].Type = getStockType(reqBody.TrackInventory)
	stock[0].RecordType = getRecordType(reqBody.TrackInventory)
	return stock
}

func setRequestStockFieldValues(newInstance *schema.Stock, reqStock *models.AddStockRequestBody, requestParams commonModel.AccountRequstParams) {
	newInstance.OwnerID = requestParams.CreatedBy
	newInstance.AccountID = requestParams.AccountID
	newInstance.Name = reqStock.Name
	newInstance.SKU = reqStock.SKU
	newInstance.Description = reqStock.Description
	newInstance.SellingPrice = reqStock.SellingPrice.Amount
	newInstance.SellingPriceCurrency = reqStock.SellingPrice.Currency
	newInstance.PurchasePrice = reqStock.PurchasePrice.Amount
	newInstance.PurchasePriceCurrency = reqStock.PurchasePrice.Currency
	newInstance.PurchaseDate = reqStock.PurchaseDate
	newInstance.ExpiryDate = reqStock.ExpiryDate
	newInstance.TrackInventory = reqStock.TrackInventory
	newInstance.StockQty = reqStock.StockQty
	newInstance.ReorderQty = reqStock.ReorderQty
	newInstance.Type = getStockType(reqStock.TrackInventory)
	newInstance.RecordType = getRecordType(reqStock.TrackInventory)
	newInstance.AsOfDate = reqStock.AsOfDate
	newInstance.AttachmentKey = reqStock.AttachmentKey
	newInstance.UnitID = reqStock.UnitID
	newInstance.CategoryID = reqStock.CategoryID
	newInstance.LocationID = reqStock.LocationID
	newInstance.SupplierID = reqStock.SupplierID
	newInstance.CreatedAt = time.Now()
	newInstance.CreatedBy = requestParams.CreatedBy
	newInstance.GroupItemID = reqStock.GroupItemID
	newInstance.SaleTaxID = reqStock.SaleTaxID
	newInstance.PurchaseTaxID = reqStock.PurchaseTaxID
	newInstance.StockUUID = uuid.New()
}

func getRecordType(inventory bool) string {
	if inventory {
		return stockConst.RecordTypeStock
	} else {
		return stockConst.RecordTypeProduct
	}
}

func getStockType(inventory bool) string {
	if inventory {
		return stockConst.TypeInventory
	} else {
		return stockConst.TypeNonInventory
	}
}

func getStockActivityRequest(stock *schema.Stock, requestParams commonModel.AccountRequstParams) *schema.StockActivity {
	var stockActivity schema.StockActivity

	currentDate := time.Now()
	stockActivity.AccountID = requestParams.AccountID
	stockActivity.OwnerID = requestParams.CreatedBy
	stockActivity.Mode = stockConst.StockActivityModeCreate
	stockActivity.OperationType = stockConst.StockActivityOperationTypeAdd
	stockActivity.StockID = stock.ID
	stockActivity.PurchaseDate = stock.PurchaseDate
	stockActivity.LocationID = stock.LocationID
	stockActivity.QuantityOnHand = 0
	stockActivity.NewQuantity = stock.StockQty
	stockActivity.AdjustedQuantity = stock.StockQty
	stockActivity.PurchasePreviousValue = decimal.NewFromFloat32(0.0)
	stockActivity.PurchasePreviousValueCurrency = stock.PurchasePriceCurrency
	stockActivity.PurchaseNewValue = stock.PurchasePrice
	stockActivity.PurchaseNewValueCurrency = stock.PurchasePriceCurrency
	stockActivity.PurchaseAdjustedValue = stock.PurchasePrice
	stockActivity.PurchaseAdjustedValueCurrency = stock.PurchasePriceCurrency
	stockActivity.SellingPreviousValue = decimal.NewFromFloat32(0.0)
	stockActivity.SellingPreviousValueCurrency = stock.SellingPriceCurrency
	stockActivity.SellingNewValue = stock.SellingPrice
	stockActivity.SellingNewValueCurrency = stock.SellingPriceCurrency
	stockActivity.SellingAdjustedValue = stock.SellingPrice
	stockActivity.SellingAdjustedValueCurrency = stock.SellingPriceCurrency
	stockActivity.AdjustedDate = &currentDate
	return &stockActivity
}
