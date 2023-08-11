package service

import (
	"pi-inventory/modules/profile/cache"
	"pi-inventory/modules/profile/repository"
)

type ProfileServiceInterface interface {
}

type profileService struct {
	profileRepository repository.ProfileRepositoryInterface
	prifileCacheRepo  cache.ProfileCacheRepositoryInterface
}

func NewProfileService(profileRepo repository.ProfileRepositoryInterface,
	prifileCacheRepo cache.ProfileCacheRepositoryInterface) *profileService {
	return &profileService{
		profileRepository: profileRepo,
		prifileCacheRepo:  prifileCacheRepo,
	}
}
