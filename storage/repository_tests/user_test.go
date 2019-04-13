package repository_tests

import (
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"testing"
)

func TestCreateDeleteUpdateUser(t *testing.T) {
	db := getDatabaseFromArgs()
	if db == nil {
		t.Error("Failed to connect to database")
	}

	user := &models.User{
		Name:     "test_user",
		Email:    "test@test.test",
		IsActive: true,
		IsAdmin:  false,
	}

	// Create user for first time
	// Should succeed
	err := db.User.Create(user, "password")
	if err != nil {
		t.Errorf("cannot create user: %s", err.Error())
	}

	user, err = db.User.FindByName("test_user")
	if err != nil {
		t.Error(Err.Wrap(&err, "Failed to get user from db"))
	}

	// Try to create duplicate user
	// Should fail
	err = db.User.Create(user, "password")
	if err == nil {
		t.Error("failed to deny creating same user again")
	}

	e := err.(*Err.Error)
	if e.Code != Err.Econflict {
		t.Error("Invalid error from creating same user again, expected EConlift, got: ", e.Code)
	}

	// Update user data
	user.Email = "test2@test2.test2"
	err = db.User.Update(user)
	if err != nil {
		t.Error("Failed to update user", err.Error())
	}

	user, err = db.User.FindByName("test_user")
	if err != nil {
		t.Error(Err.Wrap(&err, "Failed to get user from db"))
	}
	if user.Email != "test2@test2.test2" {
		t.Error("User data not updated")
	}

	// Soft delete user
	err = db.User.Delete(user)
	if err != nil {
		t.Error("Failed to delete user: ", err.Error())
	}

	// Hard delete user
	db.RemoveAllRecords()

}
