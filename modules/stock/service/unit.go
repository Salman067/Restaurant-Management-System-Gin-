package service

import (
	"fmt"
	"net/http"
	"pi-inventory/common/cache"
	"pi-inventory/common/logger"
	"pi-inventory/errors"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
	"pi-inventory/modules/stock/schema"
	"time"

	"gorm.io/gorm"
	commonModel "pi-inventory/common/models"
)

type UnitServiceInterface interface {
	CreateUnit(requestParams commonModel.AccountRequstParams, reqUnit *models.AddUnitRequestBody) (*schema.Unit, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.UnitQueryParams) ([]*models.UnitResponseBody, *commonModel.PageResponse, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.UnitResponseBody, error)
	UpdateUnit(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateUnitRequestBody) (*schema.Unit, error)
	DeleteUnit(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
}

type unitService struct {
	unitRepository  repository.UnitRepositoryInterface
	stockRepository repository.StockRepositoryInterface
	cache           cache.CacheInterface
}

func NewUnitService(unitRepo repository.UnitRepositoryInterface, stockRepository repository.StockRepositoryInterface, cache cache.CacheInterface) *unitService {
	return &unitService{
		unitRepository:  unitRepo,
		stockRepository: stockRepository,
		cache:           cache,
	}
}

func (us *unitService) CreateUnit(requestParams commonModel.AccountRequstParams, reqUnit *models.AddUnitRequestBody) (*schema.Unit, error) {
	err := reqUnit.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	_, err = us.unitRepository.FindBy(requestParams, stockConst.UnitFieldTitle, reqUnit.Title)

	if err == nil {
		logger.LogError("A unit with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "unit",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	newUnit := &schema.Unit{
		Title:     reqUnit.Title,
		OwnerID:   requestParams.CreatedBy,
		AccountID: requestParams.AccountID,
		IsUsed:    reqUnit.IsUsed,
		CreatedAt: time.Now(),
		CreatedBy: requestParams.CreatedBy,
	}
	unit, err := us.unitRepository.CreateUnit(newUnit)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "unit",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Unit, requestParams.CreatedBy, unit.ID)
	err = us.cache.Set(cacheKey, unit.Title)
	if err != nil {
		logger.LogError(err)
	}
	return unit, nil
}

func (us *unitService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.UnitQueryParams) ([]*models.UnitResponseBody, *commonModel.PageResponse, error) {
	var jsonResponse []*models.UnitResponseBody
	units, recordCount, err := us.unitRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	for _, unit := range units {
		resUnit := models.UnitResponseBody{
			ID:        unit.ID,
			OwnerID:   unit.OwnerID,
			AccountID: unit.AccountID,
			Title:     unit.Title,
			Status:    unit.Status,
			IsUsed:    unit.IsUsed,
		}
		jsonResponse = append(jsonResponse, &resUnit)
	}
	pageInfo := &commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}
	return jsonResponse, pageInfo, nil
}

func (us *unitService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.UnitResponseBody, error) {

	unit, err := us.unitRepository.FindBy(requestParams, stockConst.UnitFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "unit",
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

	cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Unit, requestParams.CreatedBy, ID)
	err = us.cache.Set(cacheKey, unit.Title)
	if err != nil {
		logger.LogError(err)
	}

	resUnit := &models.UnitResponseBody{
		ID:        unit.ID,
		OwnerID:   unit.OwnerID,
		AccountID: unit.AccountID,
		Title:     unit.Title,
		Status:    unit.Status,
		IsUsed:    unit.IsUsed,
	}
	return resUnit, nil
}

func (us *unitService) UpdateUnit(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateUnitRequestBody) (*schema.Unit, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	unit, err := us.unitRepository.FindBy(requestParams, stockConst.UnitFieldID, ID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "unit",
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

	checkedUnit := updateUnitField(reqBody, unit)
	unit, err = us.unitRepository.UpdateUnit(checkedUnit)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "unit",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Unit, requestParams.CreatedBy, ID)
	err = us.cache.Set(cacheKey, unit.Title)
	if err != nil {
		logger.LogError(err)
	}
	return unit, nil
}

func (us *unitService) DeleteUnit(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	_, err := us.unitRepository.FindBy(requestParams, stockConst.UnitFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "unit",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError}
	}
	//Check Before Delete
	stocks, err := us.stockRepository.FindBy(requestParams, stockConst.StockFieldUnitID, ID)
	if err != nil {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	} else if len(stocks) > 0 {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteFailedDueToDependency",
			TranslationParams: map[string]interface{}{
				"field": "unit",
			},
			HttpCode: http.StatusBadRequest}
	} else {
		ID, err = us.unitRepository.DeleteUnit(requestParams, ID)
		if err != nil {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.Unsuccessfull,
				TranslationKey: "deleteUnsuccessfull",
				TranslationParams: map[string]interface{}{
					"field": "unit",
				},
				HttpCode: http.StatusInternalServerError}
		}
		cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Unit, requestParams.CreatedBy, ID)
		err = us.cache.Remove(cacheKey)
		if err != nil {
			logger.LogError(err)
		}
		return ID, nil
	}
}

func updateUnitField(reqBody *models.UpdateUnitRequestBody, unit *schema.Unit) *schema.Unit {
	if reqBody.Title != "" {
		unit.Title = reqBody.Title
	}
	if reqBody.Status != "" {
		unit.Status = reqBody.Status
	}
	unit.IsUsed = reqBody.IsUsed
	unit.UpdatedAt = time.Now()
	return unit
}
