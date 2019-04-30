package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tryffel/fusio/util"
	"time"
)

type Array struct {
	Items []interface{}
}

type Pipeline struct {
	Id        string `gorm:"primary key"`
	Owner     User
	OwnerId   uint   `gorm:"not null"`
	Enabled   bool   `gorm:"not null"`
	Name      string `gorm:"not null"`
	LowerName string `gorm:"not null"`
	Info      string
	// Data stores actual json
	Data postgres.Jsonb `gorm:"not null"`
	// Blocks have array of blocks
	Blocks    Array     `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (p *Pipeline) BeforeCreate() error {
	p.Id = util.NewUuid()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Pipeline) BeforeUpdate() error {
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Pipeline) AfterFind() error {
	return nil

}
