package Err

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Emsg string

const (
	Econflict     Emsg = "conflict"  // action cannot be performed
	Einternal     Emsg = "internal"  // internal error
	Einvalid      Emsg = "invalid"   // validation failed
	Enotfound     Emsg = "not_found" // entity does not exist
	EconflictCode int  = 10
	EinternalCode int  = 20
	EnotfouncCode int  = 30
	EinvalidCode  int  = 40
)

type ErrorResolver interface {
	// EnsUserMessage returns end user friendly (and safe) error message
	EndUserMessage() string
	// EndUserCode returns E<>Code for end user
	EndUserCode() string
	// LogMessage logs error message if needed.
	LogMessage() string
	// Wrap existing error in Err.Error, and create new Err.Error if needed.
	Wrap(reason string)
	// Get cause (initial error) for error
	Cause() string
	// Implement error interface
	Error() string
}

// Application error.
type Error struct {
	// Standard code for both internal and external use
	Code Emsg

	// Message for end user
	Message string

	// Error stack for logs
	Err error
}

func (e *Error) EndUserMessage() string {
	if e.Message == "" {
		switch e.Code {
		case Econflict:
			return "Conflict"
		case Einvalid:
			return "Invalid request"
		case Einternal:
			return "Internal error"
		case Enotfound:
			return "Resource not found"
		}
	}
	return e.Message
}

func (e *Error) EndUserCode() int {
	switch e.Code {
	case Econflict:
		return EconflictCode
	case Einternal:
		return EinternalCode
	case Einvalid:
		return EinvalidCode
	case Enotfound:
		return EnotfouncCode
	default:
		return EinternalCode
	}
}

func (e *Error) LogMessage() string {
	return e.Err.Error()
}

func (e *Error) Wrap(reason string) {
	e.Err = errors.Wrap(e.Err, reason)
}

func (e *Error) Cause() string {
	if e.Err != nil {
		return errors.Cause(e.Err).Error()
	}
	return e.Error()
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func (e *Error) EndUserError() (string, error) {
	data := map[string]interface{}{}
	data["code"] = e.EndUserCode()
	data["message"] = e.EndUserMessage()

	bytes, err := json.Marshal(data)
	if err != nil {
		text := string(bytes)
		return text, nil
	}
	return "", err
}

// LogError logs all unidentified errors automatically and
// Err.Error are logged only if Error.code == Einternal
func Log(err error) {
	if err == nil {
		return
	}
	if e, ok := err.(*Error); ok {
		if e.Code == Einternal {
			logrus.Error(err)
		} else if e.Code == Econflict {
			logrus.Debug(err)
		}
		return
	}
	logrus.Error(err)
}

// Wrap existing error along with text. If passed error is already Err.Error, return same error. Otherwise create
// Err.Error and return it, Error.Code is then Einternal
func Wrap(err *error, text string) *Error {
	if e, ok := (*err).(*Error); ok {
		e.Wrap(text)
		return e
	}
	return &Error{Code: Einternal, Err: *err}
}

// GetErrCode returns Emsg if there is one. If it's missing, return Einternal
func GetErrCode(err error) Emsg {
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return Einternal
}

func InternalError(reason string, err error) error {
	e := Error{
		Code: Einternal,
	}
	if err == nil {
		e.Err = errors.New(reason)
	} else {
		e.Err = Wrap(&err, reason)
	}
	return &e
}

func UserError(reason string, err error) error {
	e := Error{
		Code:    Einvalid,
		Message: reason,
		Err:     err,
	}
	return &e
}
