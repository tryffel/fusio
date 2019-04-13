package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
)

type outputChannel struct {
	db *gorm.DB
}

func (c *outputChannel) Create(channel *models.OutputChannel) error {
	return c.db.Create(&channel).Error
}

func (c *outputChannel) FindByOwner(id uint) (*[]models.OutputChannel, error) {
	out := &[]models.OutputChannel{}
	res := c.db.Where("owner_id = ?", id)
	return out, res.Error
}

func (c *outputChannel) FindbyId(id string) (*models.OutputChannel, error) {
	out := &models.OutputChannel{}
	res := c.db.Where("id = ?", id)
	return out, res.Error
}

func (c *outputChannel) FindbyOwnerAndId(ownerId uint, id string) (*models.OutputChannel, error) {
	out := &models.OutputChannel{}
	res := c.db.Where("owner_id = ? and id = ?", ownerId, id)
	return out, res.Error
}

func (c *outputChannel) Remove(channel *models.OutputChannel) error {
	panic("implement me")
}

func (c *outputChannel) Update(channel *models.OutputChannel) error {
	panic("implement me")
}

func NewOutputChannelRepository(db *gorm.DB) repository.OutputChannel {
	return &outputChannel{db: db}
}
