package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

type AlarmFilter struct {
	Filters []Influxdb.Filter `json:"filters"`
	//Trigger float32 `json:"trigger"`
	Expression string `json:"expression"`
	Limit      int64  `json:"limit"`
}

// Format used for govaluate: mean(temp) -> mean_temp
//func (i *Input) ToId() string {
//	return fmt.Sprintf("%s_%s", i.Aggregation, i.Key)
//}

type Alarm struct {
	ID        string `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	LowerName string `gorm:"not null"`
	Info      string
	Message   string `gorm:"not null"`
	OwnerId   uint   `gorm:"not null"`
	Group     string `gorm:"not null"`
	Fired     bool   `gorm:"not null"`
	Enabled   bool   `gorm:"not null, default:'true'"`
	//Query       AlarmQuery `json:"query"`
	Filter      AlarmFilter `gorm:"type:text" json:"filter"`
	History     []AlarmHistory
	RunInterval time.Duration `gorm:"not null"`

	LastRun   time.Time
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// ALarmQuery is json formatted struct that encapsulates everything needed to build actual alarm query
type AlarmQuery struct {
	Group   string            `json:"group"`
	Filters []Influxdb.Filter `json:"inputs"`
	//Range    time.Duration     `json:"range"`
	Interval time.Duration `json:"interval"`
	// Expression is govaluate-valid expression: 'mean(temperature)>10'
	Expression string `json:"expression"`
	Limit      int64  `json:"limit"`
}

/*
func (a *AlarmQuery) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("Invalid type assertion when reading alarm from db")
	}
	b := []byte(str)
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	return nil
}

func (a AlarmQuery) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}
*/

func (f *AlarmFilter) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("AlarmFilter is not string in db (alarms.filter)")
	}
	b := []byte(str)
	err := json.Unmarshal(b, f)
	if err != nil {
		return err
	}
	return nil
}

func (f AlarmFilter) Value() (driver.Value, error) {
	j, err := json.Marshal(f)
	return j, err
}

func (a *Alarm) BeforeCreate() error {
	a.LastRun = time.Now()
	return nil
}

// BeforeCreate hook that gets called when creating new instance
func (a *Alarm) BeforeSave() (err error) {
	a.ID = util.NewUuid()
	a.LowerName = strings.ToLower(a.Name)
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

// BeforeUpdate hook that updates timestamp
func (a *Alarm) BeforeUpdate() (err error) {
	a.UpdatedAt = time.Now()
	return
}

func (a *Alarm) ToAlarmQuery() *AlarmQuery {
	q := &AlarmQuery{
		Group:      a.Group,
		Filters:    a.Filter.Filters,
		Limit:      a.Filter.Limit,
		Expression: a.Filter.Expression,
		Interval:   a.RunInterval,
	}
	return q
}

func GetAlarmById(db *gorm.DB, id string) (*Alarm, error) {
	alarm := &Alarm{}
	return alarm, db.Where("id = ?", id).First(&alarm).Error
}

// GetAlarmsToValuate Get alarms that requires valuation with given interval
// Interval defines the minimum time range for duration to wait before valuating same alarms
func GetAlarmsToValuate(db *gorm.DB, interval time.Duration) *[]Alarm {
	alarms := &[]Alarm{}
	db.Where(" enabled = True AND last_run + interval '1s'*run_interval/1000000000 < now()"+
		"AND last_run + interval '1s' * ? < now()", fmt.Sprintf("%d", int64(interval.Seconds()))).Find(alarms)
	return alarms
}

// Fire alarm: create alarmHistory
// If alarm is already fired, don't create new event
// Only one fired event can be uncleared at time. That is, one needs to first clear old fire when firing alarm again
// Timestamp: time when alarm fired
func (a *Alarm) Fire(db *gorm.DB, value float32, timestamp time.Time) error {
	if a.Fired {
		return nil
	}
	h := &AlarmHistory{
		AlarmId: a.ID,
		Value:   fmt.Sprintf("%f", value),
		Cleared: false,
		FiredAt: timestamp,
	}
	db.Where("alarm_id = ? AND clear = False").First(&h)
	if h.ID > 0 {
		return errors.New("Alarm not cleared yet. Cannot fire alarm before cleared old one!")
	}
	db.Create(h)

	a.Fired = true
	res := db.Model(*a).Update("fired", true)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Clear alarm. If alarm.fired = false, no action is taken
func (a *Alarm) Clear(db *gorm.DB, timestamp time.Time) error {
	if a.Fired == false {
		return nil
	}
	h := &AlarmHistory{}
	db.Where("alarm_id = ? AND cleared = False", a.ID).First(h)
	if h.ID < 1 {
		return errors.New(fmt.Sprintf("No alarm_history found for alarm ", a.ID))
	}

	h.Cleared = true
	h.ClearedAt = timestamp

	res := db.Save(&h)
	if res.Error != nil {
		return res.Error
	}

	a.Fired = false
	res = db.Model(*a).Update("fired", false)
	return res.Error
}

func (a *Alarm) UpdateRunTs(db *gorm.DB, timestamp time.Time) error {
	a.LastRun = timestamp
	return db.Model(*a).Update("last_run", a.LastRun).Error
}
