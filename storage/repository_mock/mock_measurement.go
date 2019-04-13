package repository_mock

import (
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/models"
	"time"
)

// Mock implementation of Measurement repository
// Currently does not do any validation and thus never returns error
// Reading values not supported yet
type MockMeasurementRepository struct {
	Measurements map[string]float64
	Metrics      map[string]float64
	StoreResults bool
}

func NewMockMeasurementRepository() *MockMeasurementRepository {
	return &MockMeasurementRepository{
		Measurements: make(map[string]float64),
		Metrics:      make(map[string]float64),
	}
}

func (m *MockMeasurementRepository) Write(device *models.Device, measurements Influxdb.Measurements) error {
	if !m.StoreResults {
		return nil
	}
	return nil
}

func (m *MockMeasurementRepository) Read(device string, group string, filters []Influxdb.Filter, from time.Time, to time.Time, n int64) (Influxdb.Batch, error) {
	panic("implement me")
}

func (m *MockMeasurementRepository) WriteMetrics(name string, value float64) error {
	panic("implement me")
}

func (m *MockMeasurementRepository) WriteMetricsBatch(batch *map[string]float64) error {
	if !m.StoreResults {
		return nil
	}
	for i, v := range *batch {
		m.Metrics[i] = v
	}
	return nil
}

func (m *MockMeasurementRepository) GetDeviceMeasurements(device string) ([]string, error) {
	panic("implement me")
}

func (m *MockMeasurementRepository) GetGroupMeasurements(group string) ([]string, error) {
	panic("implement me")
}
