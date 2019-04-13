package dtos

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/govalidator"
	"github.com/tryffel/fusio/util"
	"net/http"
	"time"
)

// Validator interface
type Validator interface {
	// Returns Validation map
	ValidationMap() *govalidator.MapData
	// Returns Validation messages to be sent to client in case of invalid data
	ValidationMessages() *govalidator.MapData
}

// Validate dto from request. If validation fails, automatically write error msg into response
// Custom validators available:
// uuid_array
func Validate(w http.ResponseWriter, r *http.Request, v Validator) error {

	rules := v.ValidationMap()
	msg := v.ValidationMessages()

	opts := govalidator.Options{
		Request:         r,
		Rules:           *rules,
		Data:            v,
		Messages:        *msg,
		RequiredDefault: false,
	}

	va := govalidator.New(opts)

	e := va.ValidateJSON()
	if len(e) > 0 {
		err := map[string]interface{}{"invalid_request_error": e}
		w.Header().Set("Content-type", "application/json")
		error := json.NewEncoder(w).Encode(err)
		if err != nil {
			logrus.Error(error)
		}

		return errors.New("invalid body")
	}
	return nil
}

func AddUuidArrayValidator() {
	govalidator.AddCustomRule("uuid_array", func(field string, rule string, message string, value interface{}) error {
		for _, v := range value.([]string) {
			if !util.IsUuid(v) {
				return errors.New("invalid uuid")
			}
		}
		return nil
	})
}

func AddIntervalValidator() {
	govalidator.AddCustomRule("interval", func(field string, rule string, message string, value interface{}) error {
		_, ok := value.(util.Interval)
		if !ok {
			return errors.New("Invalid interval")
		}
		return nil
	})

}

func AddDurationValidator() {
	govalidator.AddCustomRule("duration", func(field string, rule string, message string, value interface{}) error {
		_, ok := time.ParseDuration(value.(string))
		if ok != nil {
			return errors.New("Invalid interval")
		}
		return nil
	})

}
