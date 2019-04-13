package repository

import (
	"github.com/tryffel/fusio/storage/models"
	"time"
)

// Alarm interface
type Alarm interface {
	Create(alarm *models.Alarm) error
	Update(alarm *models.Alarm) error
	Remove(alarm *models.Alarm) error
	FindById(id string) (*models.Alarm, error)
	FindByOwner(id int) (*[]models.Alarm, error)
	FindByOwnerAndId(id string, owner int) (*models.Alarm, error)
	// Fire alarm
	Fire(alarm *models.Alarm, value float32, timestamp time.Time) error
	// Clear firing alarm
	Clear(alarm *models.Alarm, timestamp time.Time) error
	// GetAlarmsToValuate gets all alarms that should be evaluated withing defined interval
	GetAlarmsToValuate(interval time.Duration) (*[]models.Alarm, error)

	// LoadHistory loads history items for alarm
	LoadHistory(alarm *models.Alarm) error

	// GetHistorySize gets number of historical events for alarm
	GetHistorySize(alarm *models.Alarm) (int, error)

	UpdateRunTimestamp(alarm *models.Alarm, timestamp time.Time) error
}
