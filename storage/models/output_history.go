package models

import "github.com/jinzhu/gorm"

//OutputHistory represents historical data for outputs
type OutputHistory struct {
	gorm.Model
	Output   Output
	OutputId string `gorm:"not null"`
	Success  bool   `gorm:"not null"`
	Message  string
}
