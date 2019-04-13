package dtos

import (
	"github.com/thedevsaddam/govalidator"
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

type NewAlarm struct {
	Name    string `json:"name"`
	Info    string `json:"info"`
	Group   string `json:"group"`
	Message string `json:"message"`
	Enabled bool   `json:"enabled"`
	// Interval: how often to evaluate alarm: e.g. 2 min
	Interval string `json:"interval"`
	// Trigger: how many positive evaluations to count before firing alarm:
	// e.g. 5 with 2 min interval = 5 consecutive positive fires after 10 min evaluation
	Trigger int64  `json:"trigger"`
	Filter  string `json:"filter"`
}

func (n *NewAlarm) ToAlarm() (*models.Alarm, error) {

	dur, err := time.ParseDuration(n.Interval)
	if err != nil {
		return &models.Alarm{}, err
	}

	a := &models.Alarm{
		Name:        n.Name,
		Info:        n.Info,
		Group:       n.Group,
		Message:     n.Message,
		Enabled:     n.Enabled,
		RunInterval: dur,
	}
	inputs, err := Influxdb.FilterFromString(n.Filter)
	if err != nil {
		return a, err
	}

	// Simplify filter clause: mean(temp) -> mean_temp
	filter := strings.Replace(n.Filter, " ", "", -1)
	for _, v := range *inputs {
		s := v.String()
		sim := v.StringSimplified()
		filter = strings.Replace(filter, s, sim, -1)
	}

	af := models.AlarmFilter{
		Filters:    *inputs,
		Expression: filter,
		Limit:      n.Trigger,
	}
	a.Filter = af
	return a, nil
}

func (n *NewAlarm) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"name":    []string{"required"},
		"info":    []string{},
		"group":   []string{"uuid", "required"},
		"message": []string{},
		// Interval has to be between 1s-1y
		"interval": []string{"duration"},
		"trigger":  []string{"required"},
		// Expression will be should when changing to alarm
		"filter": []string{"regex:.+[+-><=].+", "required"},
	}
}

func (n *NewAlarm) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"name":     []string{"Descriptive name for alarm. Required"},
		"info":     []string{"Additional information about alarm"},
		"group":    []string{"Group that alarms get's valuated inside"},
		"message":  []string{"Message that get's sent when alarm gets fired"},
		"interval": []string{"Interval for evaluating alarm. E.g. '1m' means alarm is evaluated every 1m"},
		"trigger": []string{"Trigger for how many consecutive positive before alarming. E.g. with interval of 1m " +
			"and trigger of 10, after 10 min of positive evaluations alarm will get fired. Set to 1 to immediately " +
			"fire alarm after one positive evaluation"},
		"filter": []string{"Expression for evaluation. e.g. 'mean(temperature) - max(humidity) > 10'"},
	}
}

type Alarm struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Info     string        `json:"info"`
	Message  string        `json:"message"`
	Fired    bool          `json:"fired"`
	Enabled  bool          `json:"enabled"`
	Group    string        `json:"group"`
	Interval util.Interval `json:interval`
	Trigger  int64         `json:"trigger"`
	Filter   string        `json:"filter"`
}

func AlarmToDto(a *models.Alarm) *Alarm {
	alarm := &Alarm{
		Id:       a.ID,
		Name:     a.Name,
		Info:     a.Info,
		Message:  a.Message,
		Fired:    a.Fired,
		Enabled:  a.Enabled,
		Group:    a.Group,
		Interval: util.Interval(a.RunInterval),
	}
	return alarm
}
