package service

import (
	"context"
	"encoding/json"
	"pi-inventory/common/logger"
	"pi-inventory/modules/stock/cache"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
)

type SupplierServiceInterface interface {
	FindAll(requestParams models.RequstParams) ([]*models.SupplierResponseBody, error)
	FindByID(requestParams models.RequstParams, ID string) (*models.SupplierResponseBody, error)
}
type supplierService struct {
	supplierRepository      repository.SupplierRepositoryInterface
	supplierCacheRepository cache.SupplierCacheRepositoryInterface
}

func NewSupplierService(supplierRepo repository.SupplierRepositoryInterface, supplierCacheRepository cache.SupplierCacheRepositoryInterface) *supplierService {
	return &supplierService{supplierRepository: supplierRepo, supplierCacheRepository: supplierCacheRepository}
}

func (ss *supplierService) FindAll(requestParams models.RequstParams) ([]*models.SupplierResponseBody, error) {
	supplierList := make([]*models.SupplierResponseBody, 0)

	key := "supplier_" + requestParams.AccountSlug // Replace with your desired key pattern

	// Retrieve keys that match the pattern
	// keys, err := ss.supplierCacheRepository.GetKeys(context.Background(), keyPattern)
	// if err != nil {
	// 	logger.LogError(err)
	// }

	// if len(keys) == 0 {
	// 	return getDemoSupplierList(), nil
	// }

	// var supplierKey string
	// for _, key := range keys {
	// 	if strings.HasPrefix(key, "supplier_") {
	// 		supplierKey = key
	// 		break
	// 	}
	// }

	supplierListMap, err := ss.supplierCacheRepository.GetAll(context.Background(), key)
	if err != nil {
		logger.LogError(err)
	}

	for _, val := range supplierListMap {
		var supplierResp models.SupplierResponseBody
		err := json.Unmarshal([]byte(val), &supplierResp)
		if err != nil {
			logger.LogError("Error while unmarshal the json in tax response")
		}
		supplierList = append(supplierList, &supplierResp)
	}

	return supplierList, nil
}

func (ss *supplierService) FindByID(requestParams models.RequstParams, ID string) (*models.SupplierResponseBody, error) {

	key := "supplier_" + requestParams.AccountSlug // Replace with your desired key pattern

	// Retrieve keys that match the pattern
	// keys, err := ss.supplierCacheRepository.GetKeys(context.Background(), keyPattern)
	// if err != nil {
	// 	logger.LogError(err)
	// }

	// if len(keys) == 0 {
	// 	return demoSupplier(ID), nil
	// }

	// var supplierKey string
	// for _, key := range keys {
	// 	if strings.HasPrefix(key, "supplier_") {
	// 		supplierKey = key
	// 		break
	// 	}
	// }

	supplier, err := ss.supplierCacheRepository.GetSingleData(context.Background(), key, ID)
	if err != nil {
		logger.LogError(err)
	}

	var supplierResp models.SupplierResponseBody

	err = json.Unmarshal([]byte(supplier), &supplierResp)
	if err != nil {
		logger.LogError("Error while unmarshal the json in tax response")
	}

	return &supplierResp, nil
}

// func getDemoSupplierList() []*models.SupplierResponseBody {
// 	return db.DemoSupplierList
// }

// func demoSupplier(id string) *models.SupplierResponseBody {
// 	demoSupplierList := db.DemoSupplierList

// 	for _, val := range demoSupplierList {
// 		if val.ID.String() == id {
// 			return val
// 		}
// 	}
// 	return &models.SupplierResponseBody{
// 		ID:          uuid.New(),
// 		Title:       "Demo Title",
// 		FirstName:   "Demo First Name",
// 		LastName:    "Demo Last Name",
// 		DisplayName: "Supplier 1",
// 		Status:      "active",
// 	}
// }
