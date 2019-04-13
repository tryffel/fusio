package repository

import (
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/models"
	"time"
)

// Interface for managing measurements on both relational and time-series database
type Measurement interface {
	Write(device *models.Device, measurements Influxdb.Measurements) error
	Read(device string, group string, filters []Influxdb.Filter, from time.Time, to time.Time, n int64) (Influxdb.Batch, error)
	WriteMetrics(name string, value float64) error
	WriteMetricsBatch(batch *map[string]float64) error
	GetDeviceMeasurements(device string) ([]string, error)
	GetGroupMeasurements(group string) ([]string, error)
}
