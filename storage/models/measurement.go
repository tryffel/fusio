package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// TODO: obsolete

// Measurement single measurement for given device with either boolean, string or float32 value
type Measurement struct {
	gorm.Model
	DeviceId        string `gorm:"not null"`
	GroupId         string
	Name            string
	Unit            string
	LastMeasurement time.Time
}
