package models

import "github.com/jinzhu/gorm"

type Setting struct {
	gorm.Model
	Key   string `gorm:"not null; unique"`
	Value string `gorm:"not null"`
}
