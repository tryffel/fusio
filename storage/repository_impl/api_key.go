package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"time"
)

type ApiKeyRepository struct {
	db *gorm.DB
}

func (a *ApiKeyRepository) GetDeviceId(apiKey string) (string, error) {
	key := &models.ApiKey{}
	res := a.db.Where("key = ?", apiKey).First(&key)

	if key.DeviceId == "" {
		return "", &Err.Error{Code: Err.Enotfound, Err: errors.New("not found")}
	}
	return key.DeviceId, getDatabaseError(res.Error)
}

func (a *ApiKeyRepository) New(name string, device *models.Device, expires bool, expiration time.Time) (*models.ApiKey, error) {
	key := &models.ApiKey{
		Name:     name,
		DeviceId: device.ID,
	}
	if expires {
		key.ValidTill = expiration
	}

	res := a.db.Create(&key)
	return key, getDatabaseError(res.Error)
}

func (a *ApiKeyRepository) Delete(key *models.ApiKey) error {
	err := a.db.Delete(&key).Error
	return getDatabaseError(err)
}

func (a *ApiKeyRepository) Update(key *models.ApiKey) error {
	err := a.db.Update(*key).Error
	return getDatabaseError(err)
}

func (a *ApiKeyRepository) Get(deviceId string) (*models.ApiKey, error) {
	key := &models.ApiKey{}
	res := a.db.Where("device_id = ?", deviceId).First(&key)
	return key, getDatabaseError(res.Error)
}

func NewApiKeyRepository(db *gorm.DB) repository.ApiKey {
	return &ApiKeyRepository{db: db}
}
