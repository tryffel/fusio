package models

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/util"
	"time"
)

// ApiKey keys for devices
type ApiKey struct {
	gorm.Model
	Name      string
	Device    Device
	DeviceId  string `gorm:"not null"`
	Key       string `gorm:"not null; unique"`
	LastSeen  time.Time
	ValidTill time.Time
}

func (a *ApiKey) UpdateTimestamp() {
	a.LastSeen = time.Now()
}

// BeforeCreate hook that gets called when creating new instance
func (a *ApiKey) BeforeCreate() (err error) {
	if a.Key == "" {
		a.Key = util.RandomKey(20)
	}
	return nil
}
