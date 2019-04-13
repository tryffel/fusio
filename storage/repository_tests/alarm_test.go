package repository_tests

import (
	"github.com/tryffel/fusio/storage/models"
	"testing"
)

func TestCreateAlarm(t *testing.T) {

	db := getDatabaseFromArgs()
	if db == nil {
		t.Error("Failed to open test database")
		return
	}

	alarm := &models.Alarm{
		Name:    "test_alarm",
		Info:    "testing",
		Message: "test alarm",
		OwnerId: 1,
		Fired:   false,
	}

	db.Alarm.Create(alarm)

}
