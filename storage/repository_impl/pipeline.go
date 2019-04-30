package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
)

type pipe struct {
	db *gorm.DB
}

func NewPipelineRepository(db *gorm.DB) repository.Pipeline {
	return &pipe{db: db}
}

func (m *pipe) Create(pipeline *models.Pipeline) error {
	return getDatabaseError(m.db.Create(&pipeline).Error)
}

func (m *pipe) Update(pipeline *models.Pipeline) error {
	return getDatabaseError(m.db.Update(*pipeline).Error)
}

func (m *pipe) Remove(pipeline *models.Pipeline) error {
	panic("Not implemented")
}

func (m *pipe) FindbyId(id string) (*models.Pipeline, error) {
	p := &models.Pipeline{}
	res := m.db.Where("id = ?", id).First(&p)
	return p, getDatabaseError(res.Error)
}

func (m *pipe) FindbyOwnerAndId(ownerId uint, id string) (*models.Pipeline, error) {
	p := &models.Pipeline{}
	res := m.db.Where("owner_id = ? AND id = ?", ownerId, id).First(&p)
	return p, getDatabaseError(res.Error)
}

func (m *pipe) FindByOwner(id uint) (*[]models.Pipeline, error) {
	p := &[]models.Pipeline{}
	res := m.db.Where("owner_id = ?", id).Find(&p)
	return p, getDatabaseError(res.Error)
}
