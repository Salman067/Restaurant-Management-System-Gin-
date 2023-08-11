package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pi-inventory/common/logger"
	"pi-inventory/errors"
	"pi-inventory/modules/stock/cache"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
	"pi-inventory/modules/stock/schema"
	"time"

	"gorm.io/gorm"
	commonModel "pi-inventory/common/models"
)

type PurposeServiceInterface interface {
	CreatePurpose(requestParams commonModel.AccountRequstParams, reqPurpose *models.AddPurposeRequestBody) (*schema.Purpose, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.PurposeQueryParams) ([]*models.PurposeResponseBody, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.PurposeResponseBody, error)
	DeletePurpose(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	UpdatePurpose(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdatePurposeRequestBody) (*schema.Purpose, error)
}

type purposeService struct {
	purposeRepository      repository.PurposeRepositoryInterface
	stockActivityRepo      repository.StockActivityRepositoryInterface
	purposeCacheRepository cache.PurposeCacheRepositoryInterface
}

func NewPurposeService(purposeRepo repository.PurposeRepositoryInterface,
	stockActivityRepo repository.StockActivityRepositoryInterface,
	purposeCacheRepository cache.PurposeCacheRepositoryInterface) *purposeService {
	return &purposeService{
		purposeRepository:      purposeRepo,
		stockActivityRepo:      stockActivityRepo,
		purposeCacheRepository: purposeCacheRepository,
	}
}

func (ps *purposeService) CreatePurpose(requestParams commonModel.AccountRequstParams, reqPurpose *models.AddPurposeRequestBody) (*schema.Purpose, error) {

	err := reqPurpose.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	_, err = ps.purposeRepository.FindBy(requestParams, stockConst.PurposeFieldTitle, reqPurpose.Title)

	if err == nil {
		logger.LogError("A purpose with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "reason",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	newPurpose := &schema.Purpose{
		Title:     reqPurpose.Title,
		OwnerID:   requestParams.CreatedBy,
		AccountID: requestParams.AccountID,
		IsUsed:    reqPurpose.IsUsed,
		CreatedAt: time.Now(),
		CreatedBy: requestParams.CreatedBy,
	}
	purpose, err := ps.purposeRepository.CreatePurpose(newPurpose)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Reason",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	//Storing purpose in petty-cash cache
	key := "purpose_" + requestParams.AccountSlug
	bytePurpose, err := json.Marshal(purpose)
	if err != nil {
		logger.LogError("Purpose is not marshal while storing in petty-cash cache")
	}

	stringPurposeID := fmt.Sprint(purpose.ID)
	value := map[string]interface{}{
		stringPurposeID: string(bytePurpose),
	}

	err = ps.purposeCacheRepository.Set(context.Background(), key, value)
	if err != nil {
		logger.LogError("Purpose is not set in the petty-cash cache")
	}

	return purpose, nil
}

func (ps *purposeService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.PurposeQueryParams) ([]*models.PurposeResponseBody, error) {
	var jsonResponse []*models.PurposeResponseBody
	purposes, err := ps.purposeRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	for _, purpose := range purposes {
		resPurpose := models.PurposeResponseBody{
			ID:        purpose.ID,
			OwnerID:   purpose.OwnerID,
			AccountID: purpose.AccountID,
			Title:     purpose.Title,
			IsUsed:    purpose.IsUsed,
		}
		jsonResponse = append(jsonResponse, &resPurpose)
	}
	return jsonResponse, nil
}

func (ps *purposeService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.PurposeResponseBody, error) {
	purpose, err := ps.purposeRepository.FindBy(requestParams, stockConst.PurposeFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "purpose",
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
	resPurpose := &models.PurposeResponseBody{
		ID:        purpose.ID,
		OwnerID:   purpose.OwnerID,
		AccountID: purpose.AccountID,
		Title:     purpose.Title,
		IsUsed:    purpose.IsUsed,
	}
	return resPurpose, nil
}

func (ps *purposeService) DeletePurpose(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	_, err := ps.purposeRepository.FindBy(requestParams, stockConst.PurposeFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "purpose",
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

	//Check Before Delete
	stockActivity, err := ps.stockActivityRepo.FindBy(requestParams, stockConst.StockActivityFieldPurposeID, ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.UnKnownErr,
				TranslationKey: "unknownError",
				HttpCode:       http.StatusInternalServerError,
			}
		}
	} else if len(stockActivity) > 0 {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteFailedDueToDependency",
			TranslationParams: map[string]interface{}{
				"field": "purpose",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	ID, err = ps.purposeRepository.DeletePurpose(requestParams, ID)
	if err != nil {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Reason",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	return ID, nil
}

func (ps *purposeService) UpdatePurpose(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdatePurposeRequestBody) (*schema.Purpose, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	purpose, err := ps.purposeRepository.FindBy(requestParams, stockConst.PurposeFieldID, ID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "purpose",
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
	checkedPurpose := updatePurposeField(reqBody, purpose)
	checkedPurpose.OwnerID = requestParams.CreatedBy
	checkedPurpose.AccountID = requestParams.AccountID
	purpose, err = ps.purposeRepository.UpdatePurpose(checkedPurpose)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Reason",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	//Storing purpose in petty-cash cache
	key := "purpose_" + requestParams.AccountSlug
	bytePurpose, err := json.Marshal(purpose)
	if err != nil {
		logger.LogError("Purpose is not marshal while storing in petty-cash cache")
	}

	stringPurposeID := fmt.Sprint(purpose.ID)
	value := map[string]interface{}{
		stringPurposeID: string(bytePurpose),
	}

	err = ps.purposeCacheRepository.Set(context.Background(), key, value)
	if err != nil {
		logger.LogError("Purpose is not set in the petty-cash cache")
	}

	return purpose, nil
}

func updatePurposeField(reqBody *models.UpdatePurposeRequestBody, purpose *schema.Purpose) *schema.Purpose {
	if reqBody.Title != "" {
		purpose.Title = reqBody.Title
	}
	purpose.IsUsed = reqBody.IsUsed
	purpose.UpdatedAt = time.Now()
	return purpose
}
