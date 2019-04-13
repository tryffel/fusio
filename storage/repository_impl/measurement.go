package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"time"
)

type MeasurementRepository struct {
	db     *gorm.DB
	influx Influxdb.Client
}

func (m *MeasurementRepository) GetDeviceMeasurements(device string) ([]string, error) {
	return m.influx.GetDeviceMeasurements(device)
}

func (m *MeasurementRepository) GetGroupMeasurements(group string) ([]string, error) {
	return m.influx.GetGroupMeasurements(group)
}

func (m *MeasurementRepository) Write(device *models.Device, measurements Influxdb.Measurements) error {
	if len(device.Groups) == 0 {
		device.LoadGroups(m.db)
	}
	if len(device.Groups) == 0 {
		m.influx.Write(device.ID, []string{}, measurements)
	}
	return m.influx.Write(device.ID, device.GroupIdList(), measurements)
}

func (m *MeasurementRepository) Read(device string, group string, filters []Influxdb.Filter, from time.Time, to time.Time, n int64) (Influxdb.Batch, error) {
	return m.influx.Read(device, group, filters, from, to, n)
}

func (m *MeasurementRepository) WriteMetrics(name string, value float64) error {
	return m.influx.WriteMetrics(name, value)
}

func (m *MeasurementRepository) WriteMetricsBatch(batch *map[string]float64) error {
	return m.influx.WriteMetricsBatch(batch)
}

func NewMeasurementRepository(db *gorm.DB, influx Influxdb.Client) repository.Measurement {
	return &MeasurementRepository{
		db:     db,
		influx: influx,
	}
}
