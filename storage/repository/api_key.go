package repository

import (
	"github.com/tryffel/fusio/storage/models"
	"time"
)

type ApiKey interface {
	New(name string, device *models.Device, expires bool, expiration time.Time) (*models.ApiKey, error)
	Delete(key *models.ApiKey) error
	Update(key *models.ApiKey) error
	Get(deviceId string) (*models.ApiKey, error)
	GetDeviceId(apiKey string) (string, error)
}
