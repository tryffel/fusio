package repository_mock

import (
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
	"time"
)

type MockAlarmRepository struct {
	storage []models.Alarm
}

func (r *MockAlarmRepository) UpdateRunTimestamp(alarm *models.Alarm, timestamp time.Time) error {
	panic("implement me")
}

func (r *MockAlarmRepository) LoadHistory(alarm *models.Alarm) error {
	panic("implement me")
}

func (r *MockAlarmRepository) GetHistorySize(alarm *models.Alarm) (int, error) {
	panic("implement me")
}

func (r *MockAlarmRepository) Create(alarm *models.Alarm) error {
	if alarm.ID == "" {
		alarm.ID = util.NewUuid()
	}
	r.storage = append(r.storage, *alarm)
	return nil
}

func (r *MockAlarmRepository) Update(alarm *models.Alarm) error {
	panic("implement me")
}

func (r *MockAlarmRepository) Remove(alarm *models.Alarm) error {
	panic("implement me")
}

func (r *MockAlarmRepository) FindById(id string) (*models.Alarm, error) {
	for i, v := range r.storage {
		if v.ID == id {
			return &r.storage[i], nil
		}
	}
	return &models.Alarm{}, nil
}

func (r *MockAlarmRepository) FindByOwner(id int) (*[]models.Alarm, error) {
	alarms := make([]models.Alarm, 0)
	for _, v := range r.storage {
		if v.OwnerId == uint(id) {
			alarms = append(alarms, v)
		}
	}
	return &alarms, nil
}

func (r *MockAlarmRepository) FindByOwnerAndId(id string, owner int) (*models.Alarm, error) {
	for i, v := range r.storage {
		if v.OwnerId == uint(owner) && v.ID == id {
			return &r.storage[i], nil
		}
	}
	return &models.Alarm{}, nil
}

func (r *MockAlarmRepository) Fire(alarm *models.Alarm, value float32, timestamp time.Time) error {
	alarm.Fired = true
	return nil
}

func (r *MockAlarmRepository) Clear(alarm *models.Alarm, timestamp time.Time) error {
	alarm.Fired = false
	return nil
}

// Return all alarms for evaluation
func (r *MockAlarmRepository) GetAlarmsToValuate(interval time.Duration) (*[]models.Alarm, error) {
	return &r.storage, nil
}
