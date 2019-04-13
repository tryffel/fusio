package repository

import "github.com/tryffel/fusio/storage/models"

// OutuputOpts Options to filter outputs
type OutputOpts struct {
	OnlyEnabled bool
	OnFire      bool
	OnClear     bool
	OnError     bool
}

type Output interface {
	Create(output *models.Output) error
	Update(output *models.Output) error
	Remove(output *models.Output) error
	FindbyId(id string) (*models.Output, error)
	FindbyOwnerAndId(ownerId uint, id string) (*models.Output, error)
	FindByOwner(id uint) (*[]models.Output, error)
	FindByAlarm(alarmId string, opts OutputOpts) (*[]models.Output, error)
	// Mark Output as used. Updates outputs lastPushed timestamp as well as creates
	// OutputHistory with provided data
	MarkPushed(output *models.Output, success bool, errMsg string) error
}

type OutputChannel interface {
	Create(channel *models.OutputChannel) error
	Update(channel *models.OutputChannel) error
	Remove(channel *models.OutputChannel) error
	FindbyId(id string) (*models.OutputChannel, error)
	FindbyOwnerAndId(ownerId uint, id string) (*models.OutputChannel, error)
	FindByOwner(id uint) (*[]models.OutputChannel, error)
}
