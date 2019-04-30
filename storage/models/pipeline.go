package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

type PipelineItem struct {
	Type string
	Data interface{}
}

type PipelineBlocks struct {
	Items []PipelineItem
}

// Unmarshaling Pipeline.Blocks produces PipelineItem.Data to be map[string]interface
// And each block would have to parse config from this map.

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
	Blocks    PipelineBlocks `gorm:"-"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
}

func (p *Pipeline) BeforeCreate() error {
	p.Id = util.NewUuid()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.LowerName = strings.ToLower(p.Name)
	return p.marshalBlocks()
}

func (p *Pipeline) marshalBlocks() error {
	b, err := json.Marshal(p.Blocks)
	if err != nil {
		return err
	}
	return p.Data.Scan(b)
}

func (p *Pipeline) unmarshalBlocks() error {
	val, err := p.Data.Value()
	if err != nil {
		return &Err.Error{Code: Err.Einternal, Err: errors.New("failed to parse json data")}
	}
	err = json.Unmarshal(val.([]byte), &p.Blocks)
	return err
}

func (p *Pipeline) BeforeUpdate() error {
	p.UpdatedAt = time.Now()
	return p.marshalBlocks()
}

func (p *Pipeline) AfterFind() error {
	return p.unmarshalBlocks()
}
