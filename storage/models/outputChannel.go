package models

import (
	"github.com/tryffel/fusio/util"
	"time"
)

const (
	WebHook = "web_hook"
)

type OutputChannel struct {
	ID         string `gorm:"primary_key"`
	Owner      User
	OwnerId    uint `gorm:"not null"`
	Name       string
	OutputType string `gorm:"not null"`
	Outputs    []Output
	Data       string    `gorm:"data"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

func (o *OutputChannel) BeforeCreate() error {
	o.ID = util.NewUuid()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

func (o *OutputChannel) BeforeUpdate() error {
	o.UpdatedAt = time.Now()
	return nil
}
