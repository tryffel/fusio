package repository

import "github.com/tryffel/fusio/storage/models"

type Device interface {
	Create(device *models.Device) error
	Update(device *models.Device) error
	Delete(device *models.Device) error
	GetById(id string) (*models.Device, error)
	GetByOwnerId(id uint) (*[]models.Device, error)
	LoadGroups(device *models.Device) error
	UserHasAccess(userId uint, deviceIds []string) (bool, error)
}
