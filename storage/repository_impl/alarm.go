package repository_impl

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"time"
)

type AlarmRepository struct {
	db *gorm.DB
}

func (r *AlarmRepository) UpdateRunTimestamp(alarm *models.Alarm, timestamp time.Time) error {
	alarm.LastRun = timestamp
	res := r.db.Model(*alarm).Update("last_run", timestamp)
	return getDatabaseError(res.Error)
}

func NewAlarmRepository(db *gorm.DB) repository.Alarm {
	return &AlarmRepository{db: db}
}

func (r *AlarmRepository) Clear(alarm *models.Alarm, timestamp time.Time) error {
	// Clear alarm and update alarm history
	if !alarm.Fired {
		return nil
	}

	h := &models.AlarmHistory{}
	res := r.db.Where("alarm_id = ? AND cleared = False", alarm.ID).First(h)
	err := getDatabaseError(res.Error)
	if err == nil && h.ID > 0 {

		h.Cleared = true
		h.ClearedAt = timestamp
		res := r.db.Save(&h)
		if res.Error != nil {
			return getDatabaseError(res.Error)
		}

		alarm.Fired = false
		res = r.db.Model(*alarm).Update("fired", false)
		return getDatabaseError(res.Error)

	} else {
		return Err.Wrap(&err, fmt.Sprintf("no alarm history found for alarm %s", alarm.ID))
	}
}

func (r *AlarmRepository) GetAlarmsToValuate(interval time.Duration) (*[]models.Alarm, error) {
	return models.GetAlarmsToValuate(r.db, interval), nil
}

func (r *AlarmRepository) Create(alarm *models.Alarm) error {
	res := r.db.Create(&alarm)
	return res.Error
}

func (r *AlarmRepository) Update(alarm *models.Alarm) error {
	return r.db.Update(*alarm).Error
}

func (r *AlarmRepository) Remove(alarm *models.Alarm) error {
	// Need to handle deleting alarm history
	panic("implement me")
}

func (r *AlarmRepository) FindById(id string) (*models.Alarm, error) {
	return models.GetAlarmById(r.db, id)
}

func (r *AlarmRepository) FindByOwner(id int) (*[]models.Alarm, error) {
	alarms := &[]models.Alarm{}
	res := r.db.Where("owner_id = ?", id).Find(&alarms).Limit(100)
	return alarms, res.Error
}

func (r *AlarmRepository) FindByOwnerAndId(id string, owner int) (*models.Alarm, error) {
	alarm := &models.Alarm{}
	res := r.db.Where("id = ? AND owner_id = ?", id, owner).First(&alarm)
	return alarm, res.Error
}

func (r *AlarmRepository) Fire(alarm *models.Alarm, value float32, timestamp time.Time) error {
	return alarm.Fire(r.db, value, timestamp)
}

func (r *AlarmRepository) GetHistorySize(alarm *models.Alarm) (int, error) {
	return alarm.GetHistorySize(r.db)
}

func (r *AlarmRepository) LoadHistory(alarm *models.Alarm) error {
	return alarm.LoadHistory(r.db)
}
