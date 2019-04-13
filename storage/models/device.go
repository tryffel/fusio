package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

type DeviceType string

func (d *DeviceType) ToString() string {
	return string(*d)
}

func DeviceTypeFromString(t string) (DeviceType, error) {
	if t == "sensor" {
		return DeviceSensor, nil
	} else if t == "controller" {
		return DeviceController, nil
	} else {
		return DeviceSensor, errors.New("Invalid device type")
	}
}

const (
	DeviceSensor     DeviceType = "sensor"
	DeviceController DeviceType = "controller"
)

type Device struct {
	ID         string `gorm:"primary_key"`
	Name       string `gorm:"not null"`
	LowerName  string `gorm:"not null"`
	Info       string
	OwnerId    uint       `gorm:"not null"`
	DeviceType DeviceType `gorm:"not null"`
	Groups     []Group    `gorm:"many2many:groups_devices;"`
	CreatedAt  time.Time  `gorm:"not null"`
	UpdatedAt  time.Time  `gorm:"not null"`
}

// BeforeCreate hook that gets called when creating new instance
func (d *Device) BeforeCreate() (err error) {
	d.ID = util.NewUuid()
	d.LowerName = strings.ToLower(d.Name)
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return
}

// BeforeUpdate hook that updates timestamp
func (d *Device) BeforeUpdate() (err error) {
	d.UpdatedAt = time.Now()
	return
}

// AfterFind Fill Groups map as empty array
func (d *Device) AfterFind() (err error) {
	if d.Groups == nil {
		d.Groups = []Group{}
	}
	return nil
}

func (d *Device) Exists(db *gorm.DB) bool {
	device := &Device{}
	db.Where("id = ?", d.ID).First(&device)
	if device.ID != "" {
		return true
	}
	return false
}

// Load groups loads groups associated with device
func (d *Device) LoadGroups(db *gorm.DB) error {
	groups := &[]Group{}
	res := db.Model(&d).Association("Groups").Find(&groups)
	if res.Error != nil {
		return res.Error
	}
	d.Groups = *groups
	return res.Error
}

func (d *Device) GroupIdList() []string {
	var list []string
	for _, group := range d.Groups {
		list = append(list, group.ID)
	}
	return list
}
