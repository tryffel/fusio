package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type AlarmHistory struct {
	gorm.Model
	AlarmId   string `gorm:"not null"`
	Value     string `gorm:"not null"`
	Cleared   bool   `gorm:"not null"`
	FiredAt   time.Time
	ClearedAt time.Time
}

// LoadHistory loads first 20 history items
func (a *Alarm) LoadHistory(db *gorm.DB) error {
	h := &[]AlarmHistory{}
	res := db.Where("alarm_id = ?", a.ID).Order("fired_at desc").Limit(20).Find(&h)
	a.History = *h
	return res.Error
}

func (a *Alarm) GetHistorySize(db *gorm.DB) (int, error) {
	var c int
	res := db.Where("alarm_id = ?", a.ID).Find(&[]AlarmHistory{}).Count(&c)
	return c, res.Error
}
