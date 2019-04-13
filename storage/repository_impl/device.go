package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
)

type Device struct {
	db *gorm.DB
}

func (d *Device) UserHasAccess(userId uint, deviceId []string) (bool, error) {
	devices := &[]models.Device{}
	count := 0
	res := d.db.Where("owner_id = ? AND id IN (?)", userId, deviceId).Find(devices).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	if len(deviceId) == count {
		return true, nil
	}
	return false, nil
}

func (d *Device) LoadGroups(device *models.Device) error {
	return device.LoadGroups(d.db)
}

func (d *Device) Create(device *models.Device) error {
	return d.db.Create(&device).Error
}

func (d *Device) Update(device *models.Device) error {
	return d.db.Update(*device).Error
}

func (d *Device) Delete(device *models.Device) error {
	return d.db.Delete(&device).Error
}

func (d *Device) GetById(id string) (*models.Device, error) {
	device := &models.Device{}
	res := d.db.Where("id = ?", id).First(&device)
	return device, res.Error
}

func (d *Device) GetByOwnerId(id uint) (*[]models.Device, error) {
	devices := &[]models.Device{}
	res := d.db.Where("owner_id = ?", id).Find(&devices).Limit(100)
	return devices, res.Error
}

func NewDeviceRepository(db *gorm.DB) repository.Device {
	return &Device{db: db}
}
