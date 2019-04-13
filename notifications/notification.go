package notifications

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/err"
)

type Notifier interface {
	Notify(data string) error
}

// Parse type and return notifier
func GetNotifier(NotifierType string, config string) (Notifier, error) {
	switch NotifierType {
	case "webhook":
		w := &WebHook{}
		err := json.Unmarshal([]byte(config), w)
		return w, err
	case "matrix":
		m := &Matrix{}
		err := json.Unmarshal([]byte(config), m)
		return m, err
	default:
		return nil, &Err.Error{Code: Err.Econflict, Message: "invalid output type", Err: errors.New("invalid output type")}
	}
}

// ValidateChannel validates output channel data and return message for end user if invalid
func ValidateChannel(nType string, data string) error {
	switch nType {
	case "webhook":
		hook := &WebHook{}
		err := json.Unmarshal([]byte(data), hook)
		if err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	case "matrix":
		m := &Matrix{}
		err := json.Unmarshal([]byte(data), m)
		return err
	default:
		return &Err.Error{Code: Err.Econflict, Message: "invalid output type", Err: errors.New("invalid output type")}
	}
}
