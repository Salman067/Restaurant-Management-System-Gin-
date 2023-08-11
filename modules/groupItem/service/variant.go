package service

import (
	"net/http"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	groupItemConst "pi-inventory/modules/groupItem/consts"
	"pi-inventory/modules/groupItem/models"
	"pi-inventory/modules/groupItem/repository"
	"pi-inventory/modules/groupItem/schema"
	"time"

	"gorm.io/gorm"
)

type VariantServiceInterface interface {
	CreateVariant(requestParams commonModel.AccountRequstParams, reqVariant *models.AddVariantRequestBody) (*schema.Variant, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.VariantQueryParams) ([]*models.VariantResponseBody, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.VariantResponseBody, error)
	DeleteVariant(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	UpdateVariant(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateVariantRequestBody) (*schema.Variant, error)
}

type variantService struct {
	variantRepository repository.VariantRepositoryInterface
}

func NewVariantService(variantRepo repository.VariantRepositoryInterface) *variantService {
	return &variantService{variantRepository: variantRepo}
}

func (vs *variantService) CreateVariant(requestParams commonModel.AccountRequstParams, reqVariant *models.AddVariantRequestBody) (*schema.Variant, error) {

	err := reqVariant.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	_, err = vs.variantRepository.FindBy(requestParams, groupItemConst.VariantFieldTitle, reqVariant.Title)

	if err == nil {
		logger.LogError("A Vatriant with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "attribute",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	newVariant := &schema.Variant{
		Title:     reqVariant.Title,
		OwnerID:   requestParams.CreatedBy,
		AccountID: requestParams.AccountID,
		CreatedAt: time.Now(),
		CreatedBy: requestParams.CreatedBy,
	}
	variant, err := vs.variantRepository.CreateVariant(newVariant)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Attribute",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return variant, nil
}

func (vs *variantService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.VariantQueryParams) ([]*models.VariantResponseBody, error) {
	var jsonResponse []*models.VariantResponseBody
	variants, err := vs.variantRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	for _, variant := range variants {
		resVariant := models.VariantResponseBody{
			ID:        variant.ID,
			OwnerID:   variant.OwnerID,
			AccountID: variant.AccountID,
			Title:     variant.Title,
		}
		jsonResponse = append(jsonResponse, &resVariant)
	}
	return jsonResponse, nil
}

func (vs *variantService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.VariantResponseBody, error) {
	variant, err := vs.variantRepository.FindBy(requestParams, groupItemConst.VariantFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "attribute",
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
	resVariant := &models.VariantResponseBody{
		ID:        variant.ID,
		OwnerID:   variant.OwnerID,
		AccountID: variant.AccountID,
		Title:     variant.Title,
	}
	return resVariant, nil
}

func (vs *variantService) DeleteVariant(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	_, err := vs.variantRepository.FindBy(requestParams, groupItemConst.VariantFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "attribute",
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
	ID, err = vs.variantRepository.DeleteVariant(requestParams, ID)
	if err != nil {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Attribute",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return ID, nil
}

func (vs *variantService) UpdateVariant(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateVariantRequestBody) (*schema.Variant, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	variant, err := vs.variantRepository.FindBy(requestParams, groupItemConst.VariantFieldID, ID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "attribute",
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
	checkedVariant := updateVariantField(reqBody, variant)
	variant, err = vs.variantRepository.UpdateVariant(checkedVariant)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Attribute",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return variant, nil
}

func updateVariantField(reqBody *models.UpdateVariantRequestBody, variant *schema.Variant) *schema.Variant {
	if reqBody.Title != "" {
		variant.Title = reqBody.Title
	}
	variant.UpdatedAt = time.Now()
	return variant
}
