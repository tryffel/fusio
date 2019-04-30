package repository

import "github.com/tryffel/fusio/storage/models"

type Pipeline interface {
	Create(pipeline *models.Pipeline) error
	Update(pipeline *models.Pipeline) error
	Remove(pipeline *models.Pipeline) error
	FindbyId(id string) (*models.Pipeline, error)
	FindbyOwnerAndId(ownerId uint, id string) (*models.Pipeline, error)
	FindByOwner(id uint) (*[]models.Pipeline, error)
}
