package models

import (
	"errors"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/util"
	"time"
)

const (
	BlockTypeInput  = "input"
	BlockTypeOutput = "output"
	BlockTypeFunc   = "func"
)

type PipelineBlock struct {
	Id          string `gorm:"primary key"`
	Owner       User
	OwnerId     uint   `gorm:"not null"`
	BlockType   string `gorm:"not null"`
	Pipeline    Pipeline
	PipelineId  uint           `gorm:"not null"`
	Name        string         `gorm:"not null"`
	BlockModel  string         `gorm:"not null"`
	Data        postgres.Jsonb `gorm:"not null"`
	NextBlockId string
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (p *PipelineBlock) BeforeCreate() error {
	if p.BlockType != BlockTypeInput && p.BlockType != BlockTypeFunc && p.BlockType != BlockTypeOutput {
		return &Err.Error{Code: Err.Einvalid, Err: errors.New("invalid block type")}
	}

	p.Id = util.NewUuid()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *PipelineBlock) BeforeUpdate() error {
	p.UpdatedAt = time.Now()
	return nil
}
