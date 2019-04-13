package models

import (
	"github.com/tryffel/fusio/util"
	"time"
)

type Output struct {
	ID              string `gorm:"primary_key"`
	Owner           User
	OwnerId         uint `gorm:"not null"`
	Alarm           Alarm
	AlarmId         string `gorm:"not null"`
	Name            string
	FireTemplate    string
	ClearTemplate   string
	ErrorTemplate   string
	OutputChannel   OutputChannel
	OutputChannelId string
	Repeat          time.Duration
	LastPushed      time.Time
	Enabled         bool      `gorm:"not null; default:'true'"`
	OnFire          bool      `gorm:"not null; default:'true'"`
	OnClear         bool      `gorm:"not null; default:'true'"`
	OnError         bool      `gorm:"not null; default:'true'"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (o *Output) BeforeCreate() error {
	o.ID = util.NewUuid()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Output) BeforeUpdate() error {
	o.UpdatedAt = time.Now()
	return nil
}
