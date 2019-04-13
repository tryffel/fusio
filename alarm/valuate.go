package alarm

import (
	"fmt"
	"github.com/knetic/govaluate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/metrics"
	"github.com/tryffel/fusio/notifications"
	"github.com/tryffel/fusio/storage"
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"time"
)

type Parameters map[string]interface{}
type OutputType string

const (
	Fire  OutputType = "fire"
	Clear OutputType = "clear"
	Error OutputType = "error"
)

func (p *Parameters) Get(name string) (interface{}, error) {
	if (*p)[name] == nil {
		return nil, errors.New("Not found")
	}
	return (*p)[name], nil
}

// RunAlarms checks alarms and fires / clears them if needed
func RunAlarms(alarms []models.Alarm, store storage.Store, metrics metrics.Metrics) {
	for _, v := range alarms {
		status, err, measurement := Valuate(*v.ToAlarmQuery(), store.Measurement, v.RunInterval)
		if err != nil {
			if e, ok := err.(*Err.Error); ok {
				if e.Cause() != "Alarm has not enough measurement points" {
					logrus.Errorf("Failure running alarms: %s", err.Error())
				}
			} else {
				logrus.Errorf("Failure running alarms: %s", err.Error())
			}
		}

		if err == nil {
			var val float64
			for _, v := range *measurement {
				val = v
				break
			}
			if v.Fired == false && status == true {
				logrus.Debug("Alarm ", v.ID, ", ", v.Name, " fired!")
				err = store.Alarm.Fire(&v, float32(val), time.Now())
				Err.Log(err)
				err = pushOutputs(store, &v, measurement, Fire)
				Err.Log(err)
				metrics.CounterIncrease("alarm_notification_fired", 1)
				metrics.CounterIncrease("alarm_notification", 1)
			}
			if v.Fired == true && status == false {
				logrus.Debug("Alarm ", v.ID, ", ", v.Name, " cleared!")
				err = store.Alarm.Clear(&v, time.Now())
				Err.Log(err)
				err = pushOutputs(store, &v, measurement, Clear)
				Err.Log(err)
				metrics.CounterIncrease("alarm_notification_cleared", 1)
				metrics.CounterIncrease("alarm_notification", 1)
			}
		} else {
			// Push errors
			if e, ok := err.(*Err.Error); ok {
				if e.Code == Err.Einternal {
					Err.Log(e)
				} else {
					err = pushOutputs(store, &v, measurement, Error)
					Err.Log(err)
					metrics.CounterIncrease("alarm_notification_error", 1)
					metrics.CounterIncrease("alarm_notification", 1)
				}
			}
		}
		err = store.Alarm.UpdateRunTimestamp(&v, time.Now())
		Err.Log(err)
	}
}

// Valuate evaluates single alarm and returns true if fired
func Valuate(alarmQuery models.AlarmQuery, i repository.Measurement, runInterval time.Duration) (bool, error, *map[string]float64) {
	meas, err := i.Read("", alarmQuery.Group, alarmQuery.Filters, time.Now().Add(-time.Duration(alarmQuery.Limit)*alarmQuery.Interval), time.Now(), alarmQuery.Limit)
	if err != nil {
		Err.Log(err)
		return false, err, &map[string]float64{}
	}

	if meas == nil {
		return false, nil, &map[string]float64{}
	}
	if len(meas) == 0 {
		logrus.Debug("No measurements for alarm")
		return false, nil, &map[string]float64{}
	}

	status, err, measurements := ValuateSeries(&alarmQuery, meas)
	return status, err, measurements
}

// ValuateSeries valuates series of measurements. In each point, evaluation must be true in order to return true
func ValuateSeries(query *models.AlarmQuery, batch Influxdb.Batch) (bool, error, *map[string]float64) {
	out := make(map[string]float64)
	exp, err := govaluate.NewEvaluableExpression(query.Expression)
	if err != nil {
		return false, err, &out
	}
	params := make(map[string]interface{}, len(query.Filters))
	var ts int64 = 0

	// Check batch has enough measurement points
	for _, v := range batch {
		if len(v) < int(query.Limit) {
			e := &Err.Error{Code: Err.Econflict, Err: errors.New("Alarm has not enough measurement points")}
			return false, e, &out
		}
	}
	// Evaluate
	for ts = 0; ts < query.Limit; ts++ {
		for i, v := range batch {
			params[i] = v[ts].Value
		}
		res, err := exp.Evaluate(params)
		if err != nil {
			e := Err.Wrap(&err, "Failed to evaluate alarm state")
			e.Code = Err.Einvalid
			return false, e, &out
		}
		if res == false {
			return false, nil, &out
		}
		// Last round, fill output data
		if ts == (query.Limit - 1) {
			for i, v := range batch {
				out[i] = float64(v[ts].Value)
			}
		}
	}
	return true, nil, &out
}

func pushOutputs(store storage.Store, alarm *models.Alarm, measurements *map[string]float64, outType OutputType) error {
	opts := repository.OutputOpts{}
	opts.OnlyEnabled = true
	switch outType {
	case Fire:
		opts.OnFire = true
	case Clear:
		opts.OnClear = true
	case Error:
		opts.OnError = true
	}

	outputs, err := store.Output.FindByAlarm(alarm.ID, opts)
	logrus.Debugf("Pushing %d outputs", len(*outputs))

	if err != nil {
		logrus.Errorf("Failed to retrieve alarm outputs for alarm %s: %s", alarm.ID, err.Error())
		return err
	}

	// Construct value string
	value := ""

	for i, v := range *measurements {
		value = fmt.Sprintf("%s %s=%.2f", value, i, v)
	}

	n := notifications.Notification{
		AlarmName:     alarm.Name,
		AlarmMsg:      alarm.Message,
		AlarmId:       alarm.ID,
		Timestamp:     time.Now().Format(time.Kitchen),
		TimestampUnix: fmt.Sprintf("%d", time.Now().Second()),
		Value:         value,
		Title:         "",
		GroupId:       alarm.Group,
		GroupName:     "",
		Error:         "unknown error",
	}

	for _, out := range *outputs {
		var text string
		var err error
		switch outType {
		case Fire:
			text, err = n.Parse(out.FireTemplate)
		case Clear:
			text, err = n.Parse(out.ClearTemplate)
		case Error:
			text, err = n.Parse(out.ErrorTemplate)
		}
		if err != nil {
			if e, ok := err.(*Err.Error); ok {
				e.Wrap("Cannot push outputs")
				Err.Log(e)
			} else {
				logrus.Error("Error parsing templates: ", err.Error())
			}

		} else {
			notifier, err := notifications.GetNotifier(out.OutputChannel.OutputType, out.OutputChannel.Data)
			if err != nil {
				e := Err.Wrap(&err, "Failed to get output implementation")
				Err.Log(e)
			} else {
				success := notifier.Notify(text)
				out.LastPushed = time.Now()
				err = store.Output.Update(&out)
				if err != nil {
					e := Err.Wrap(&err, "Failed to update outputs")
					Err.Log(e)
				}
				if success != nil {
					err = store.Output.MarkPushed(&out, false, success.Error())
				} else {
					err = store.Output.MarkPushed(&out, true, "")
				}

				if err != nil {
					e := Err.Wrap(&err, "Failed to mark output push")
					Err.Log(e)
				}
			}
		}
	}

	if err != nil {
		e := Err.Wrap(&err, "Failed to push outputs")
		Err.Log(e)
	}
	return nil
}
