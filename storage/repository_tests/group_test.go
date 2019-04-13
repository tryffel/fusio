package repository_tests

import (
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"testing"
)

func TestCreateUpdateDeleteGroup(t *testing.T) {
	db := getDatabaseFromArgs()
	if db == nil {
		t.Error("Failed to connect to database")
	}
	//db.LogSql(true)

	user := db.createUser()

	group := &models.Group{
		Name:    "test group",
		OwnerId: user.ID,
	}

	// Should succeed
	err := db.Group.Create(group)
	if err != nil {
		t.Errorf("Cannot create group: %s", err.Error())
	}

	group, err = db.Group.FindByOwnerAndId(user.ID, group.ID)
	if err != nil {
		t.Errorf("Failed to query group: %s", err.Error())
	}

	// Try to create duplicate
	err = db.Group.Create(group)
	if err == nil {
		t.Error("Failed to deny duplicate entry when creating group")
	}
	e := err.(*Err.Error)
	if e.Code != Err.Econflict {
		t.Error("Invalid error from creating same group again, expected 'EConflict', got: ", e.Code)
	}

	group.Info = "testing groups"
	err = db.Group.Update(group)
	if err != nil {
		t.Error("Failed to update group, ", err.Error())
	}

	group, err = db.Group.FindById(group.ID)
	if err != nil {
		t.Error("Failed to query group: ", err.Error())
	}
	if group.Info != "testing groups" {
		t.Error("Queried group doesn't match expected")
	}

	err = db.Group.Remove(group)
	if err != nil {
		t.Error("Failed to delete group: ", err.Error())
	}

	db.RemoveAllRecords()

}

func TestCreateQueryGroupDevices(t *testing.T) {
	db := getDatabaseFromArgs()
	if db == nil {
		t.Error("Failed to connect to database")
		return
	}
	//db.LogSql(true)

	user := db.createUser()
	group1 := db.createGroup(user.ID, "test_group1")
	group2 := db.createGroup(user.ID, "test_group2")
	group3 := db.createGroup(user.ID, "test_group3")

	d1 := &models.Device{
		Name:       "test_device1",
		OwnerId:    user.ID,
		DeviceType: models.DeviceSensor,
	}

	d2 := &models.Device{
		Name:       "test_device2",
		OwnerId:    user.ID,
		DeviceType: models.DeviceSensor,
	}

	d3 := &models.Device{
		Name:       "test_device3",
		OwnerId:    user.ID,
		DeviceType: models.DeviceSensor,
	}

	err := db.Device.Create(d1)
	if err != nil {
		t.Error("Failed to create new devices: ", err.Error())

	}
	err = db.Device.Create(d2)
	if err != nil {
		t.Error("Failed to create new devices: ", err.Error())
	}

	err = db.Device.Create(d3)

	if err != nil {
		t.Error("Failed to create new devices: ", err.Error())
	}

	_ = db.Group.AddDevices(group1, []string{d1.ID, d2.ID})
	_ = db.Group.AddDevices(group2, []string{d1.ID})

	devices, err := db.Group.GetDevices(user.ID, group1.ID)
	if err != nil {
		t.Error("Failed to query devices")
	}
	if len(*devices) != 2 {
		t.Error("Queried groups devices don't match")
	}

	devices, err = db.Group.GetDevices(user.ID, group2.ID)
	if err != nil {
		t.Error("Failed to query devices")
	}
	if len(*devices) != 1 {
		t.Error("Queried groups devices don't match")
	}

	devices, err = db.Group.GetDevices(user.ID, group3.ID)
	if err != nil {
		t.Error("Failed to query devices")
	}
	if len(*devices) != 0 {
		t.Error("Queried groups devices don't match")
	}

	db.RemoveAllRecords()
}
