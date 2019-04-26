package dtos

import (
	"testing"
)

func TestAlarmDtoToAlarm(t *testing.T) {
	dto := NewAlarm{
		Name:     "test",
		Info:     "information",
		Group:    "abcd-1234",
		Message:  "test message",
		Enabled:  true,
		Interval: "60s",
		Trigger:  10,
		Filter:   "mean(temperature) - derivative(max(temperature),10) > 10",
	}

	alarm, err := dto.ToAlarm()
	if err != nil {
		t.Errorf("Error creating alarm from dto: %s", err.Error())
	}

	if alarm.Name != dto.Name || alarm.Info != dto.Info || alarm.Message != dto.Message {
		t.Error("Alarm doesn't match equivalent dto")
	}

	query := alarm.ToAlarmQuery()

	if query.Group != dto.Group {
		t.Error("Alarm querys group doesn't match dto")
	}

	if query.Expression != "mean_temperature-derivative_max_temperature>10" {
		t.Errorf("alarm query doesn't match expected: %s", query.Expression)
	}
}
