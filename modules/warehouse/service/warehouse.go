package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	stockConst "pi-inventory/modules/stock/consts"
	stockModuleRepository "pi-inventory/modules/stock/repository"
	"pi-inventory/modules/warehouse/cache"
	warehouseConst "pi-inventory/modules/warehouse/consts"
	"pi-inventory/modules/warehouse/models"
	"pi-inventory/modules/warehouse/repository"
	"pi-inventory/modules/warehouse/schema"
	"time"

	"gorm.io/gorm"
)

type WarehouseServiceInterface interface {
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.WarehouseQueryParams) ([]*models.WarehouseResponseBody, *commonModel.PageResponse, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.WarehouseResponseBody, error)
	CreateWarehouse(requestParams commonModel.AccountRequstParams, reqWarehouse *models.AddWarehouseRequestBody) (*schema.Warehouse, error)
	UpdateWarehouse(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateWarehouseRequestBody) (*schema.Warehouse, error)
	DeleteWarehouse(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
}

type warehouseService struct {
	WarehouseRepository      repository.WarehouseRepositoryInterface
	stockRepo                stockModuleRepository.StockRepositoryInterface
	stockActivityRepo        stockModuleRepository.StockActivityRepositoryInterface
	warehouseCacheRepository cache.WarehouseCacheRepositoryInterface
}

func NewWarehouseService(warehouseRepo repository.WarehouseRepositoryInterface,
	stockRepo stockModuleRepository.StockRepositoryInterface,
	stockActivityRepo stockModuleRepository.StockActivityRepositoryInterface,
	warehouseCacheRepository cache.WarehouseCacheRepositoryInterface) *warehouseService {
	return &warehouseService{
		WarehouseRepository:      warehouseRepo,
		stockRepo:                stockRepo,
		stockActivityRepo:        stockActivityRepo,
		warehouseCacheRepository: warehouseCacheRepository,
	}
}

func (ws *warehouseService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.WarehouseQueryParams) ([]*models.WarehouseResponseBody, *commonModel.PageResponse, error) {
	var jsonResponse []*models.WarehouseResponseBody
	warehouses, recordCount, err := ws.WarehouseRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError}
	}
	for _, warehouse := range warehouses {
		resWarehouse := models.WarehouseResponseBody{
			ID:        warehouse.ID,
			OwnerID:   warehouse.OwnerID,
			AccountID: warehouse.AccountID,
			Title:     warehouse.Title,
			Address:   warehouse.Address,
			IsUsed:    warehouse.IsUsed,
		}
		jsonResponse = append(jsonResponse, &resWarehouse)
	}
	pageInfo := commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}
	return jsonResponse, &pageInfo, nil
}

func (ws *warehouseService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.WarehouseResponseBody, error) {
	warehouse, err := ws.WarehouseRepository.FindBy(requestParams, warehouseConst.WarehouseFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "location",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	resWarehouse := &models.WarehouseResponseBody{
		ID:        warehouse.ID,
		OwnerID:   warehouse.OwnerID,
		AccountID: warehouse.AccountID,
		Title:     warehouse.Title,
		Address:   warehouse.Address,
		IsUsed:    warehouse.IsUsed,
	}
	return resWarehouse, nil
}

func (ws *warehouseService) CreateWarehouse(requestParams commonModel.AccountRequstParams, reqWarehouse *models.AddWarehouseRequestBody) (*schema.Warehouse, error) {
	err := reqWarehouse.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	_, err = ws.WarehouseRepository.FindBy(requestParams, warehouseConst.WarehouseFieldTitle, reqWarehouse.Title)

	if err == nil {
		logger.LogError("A location with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "location",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	newWarehouse := &schema.Warehouse{
		Title:     reqWarehouse.Title,
		Address:   reqWarehouse.Address,
		OwnerID:   requestParams.CreatedBy,
		AccountID: requestParams.AccountID,
		IsUsed:    reqWarehouse.IsUsed,
		CreatedAt: time.Now(),
		CreatedBy: requestParams.CreatedBy,
	}
	warehouse, err := ws.WarehouseRepository.Create(newWarehouse)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Location",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	//Storing purpose in petty-cash cache
	key := "warehouse_"
	byteWarehouse, err := json.Marshal(warehouse)
	if err != nil {
		logger.LogError("Purpose is not marshal while storing in petty-cash cache")
	}

	stringWarehouseID := fmt.Sprint(warehouse.ID)
	value := map[string]interface{}{
		stringWarehouseID: string(byteWarehouse),
	}

	err = ws.warehouseCacheRepository.Set(context.Background(), key, value)
	if err != nil {
		logger.LogError("Purpose is not set in the petty-cash cache")
	}

	return warehouse, nil
}

func (ws *warehouseService) UpdateWarehouse(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateWarehouseRequestBody) (*schema.Warehouse, error) {

	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	warehouse, err := ws.WarehouseRepository.FindBy(requestParams, warehouseConst.WarehouseFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "location",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}

	if reqBody.Title != "" {
		warehouse.Title = reqBody.Title
	}
	if reqBody.Address != "" {
		warehouse.Address = reqBody.Address
	}
	if reqBody.AccountID != 0 {
		warehouse.AccountID = reqBody.AccountID
	}
	warehouse.IsUsed = reqBody.IsUsed
	warehouse.UpdatedAt = time.Now()
	warehouse.OwnerID = requestParams.CreatedBy
	warehouse.UpdatedBy = requestParams.CreatedBy

	warehouse, err = ws.WarehouseRepository.Update(warehouse)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Location",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	//Storing purpose in petty-cash cache
	key := "warehouse_"
	byteWarehouse, err := json.Marshal(warehouse)
	if err != nil {
		logger.LogError("Purpose is not marshal while storing in petty-cash cache")
	}

	stringWarehouseID := fmt.Sprint(warehouse.ID)
	value := map[string]interface{}{
		stringWarehouseID: string(byteWarehouse),
	}

	err = ws.warehouseCacheRepository.Set(context.Background(), key, value)
	if err != nil {
		logger.LogError("Purpose is not set in the petty-cash cache")
	}

	return warehouse, nil
}

func (ws *warehouseService) DeleteWarehouse(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	_, err := ws.WarehouseRepository.FindBy(requestParams, warehouseConst.WarehouseFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "location",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}

	stocks, stockErr := ws.stockRepo.FindBy(requestParams, stockConst.StockFieldLocationID, ID)
	_, StockActivityErr := ws.stockActivityRepo.FindBy(requestParams, stockConst.StockActivityFieldLocationID, ID)

	if stockErr != nil || StockActivityErr != nil && StockActivityErr != gorm.ErrRecordNotFound {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	} else if len(stocks) > 0 || StockActivityErr == nil {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteFailedDueToDependency",
			TranslationParams: map[string]interface{}{
				"field": "location",
			},
			HttpCode: http.StatusBadRequest,
		}
	} else {
		ID, err = ws.WarehouseRepository.Delete(requestParams, ID)
		if err != nil {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.Unsuccessfull,
				TranslationKey: "deleteUnsuccessfull",
				TranslationParams: map[string]interface{}{
					"field": "Location",
				},
				HttpCode: http.StatusInternalServerError,
			}
		}
		return ID, nil
	}

}
