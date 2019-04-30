package repository_impl

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/repository"
	"strings"
)

var errNotFound = "not found"
var errForbidden = "forbidden"
var errAlreadyExists = "already exists"
var errInternalError = "internal error"

type Error struct {
	dbType string
}

func NewErrors(dbType string) repository.Errors {
	return &Error{dbType: dbType}
}

func (e *Error) GetUserFriendlyError(err error, resource string) error {

	if e.dbType != "postgres" {
		return errors.New(errInternalError)
	}

	if is, friendly := e.IsPostgresError(err, resource); is {
		return friendly
	}

	if is, friendly := e.IsGormError(err, resource); is {
		return friendly
	}

	// If error is unknown, log it.
	logrus.Error(err)
	return errors.New(errInternalError)
}

// IsPostgresError returns true if postgres error and returns friendly error
func (e *Error) IsPostgresError(err error, resource string) (bool, error) {
	pError, ok := err.(*pq.Error)
	if !ok {
		return false, nil
	}

	if pError.Code == "23505" {
		return true, errors.New(fmt.Sprintf("'%s' %s", resource, errAlreadyExists))
	}

	logrus.Error(err)
	return true, errors.New(errInternalError)
}

func (e *Error) IsGormError(err error, resource string) (bool, error) {
	if gorm.IsRecordNotFoundError(err) || err.Error() == "record not found" {
		return true, errors.New(fmt.Sprintf("%s not found", resource))
	}
	if err == gorm.ErrInvalidSQL {
		//return true, err.n
	}
	return false, nil

}

// getGormError returns boolean true if error is produced by gorm, and maps error message to error
// If error is not from gorm, return false and nil
func getGormError(err error) (bool, error) {
	if gorm.IsRecordNotFoundError(err) || err.Error() == "record not found" {
		return true, &Err.Error{Code: Err.Enotfound, Message: "Not found", Err: errors.New("not found")}
	}
	return false, nil
}

func getPostgresError(err error) (bool, error) {
	pError, ok := err.(*pq.Error)
	if !ok {
		return false, nil
	}

	// Unique violation
	if pError.Code == "23505" {
		return true, &Err.Error{Code: Err.Econflict, Message: "Already exists", Err: errors.New("resource exists")}
	}
	return true, &Err.Error{Code: Err.Einternal, Err: errors.Wrapf(err, "Postgresl error")}
}

// Catch SQL error, always resulting in internal error
func getSqlError(err error) (bool, error) {
	if strings.Contains(err.Error(), "sql:") {
		return true, &Err.Error{Code: Err.Einternal, Err: errors.Wrap(err, "General SQL error")}
	}
	return false, nil
}

// getDatabaseError tries to first map error to gorm and postgres and only then fills error with error data
// If error = nil, return nil
func getDatabaseError(e error) error {
	if e == nil {
		return nil
	}
	g, err := getGormError(e)
	if g {
		return err
	}
	p, err := getPostgresError(e)
	if p {
		return err
	}

	sql, err := getSqlError(e)
	if sql {
		return err
	}

	// TODO: get influxdb error

	return Err.Wrap(&e, "unknown storage error")
}
