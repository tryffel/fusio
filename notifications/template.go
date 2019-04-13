package notifications

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/err"
	"text/template"
)

// All possible template fields a notification may possibly have
type Notification struct {
	AlarmName     string
	AlarmMsg      string
	AlarmId       string
	Timestamp     string
	TimestampUnix string
	Title         string
	Text          string
	Value         string
	GroupId       string
	GroupName     string
	Error         string
}

func (n *Notification) Parse(tmpl string) (string, error) {
	out, err := n.doParse(tmpl)
	if err != nil || out == "" {
		return "", &Err.Error{Code: Err.Econflict, Message: "Failed to parse template", Err: err}
	}
	return out, nil
}

func (n *Notification) doParse(tmpl string) (string, error) {
	defer recoverInvalidTemplate()
	out := bytes.NewBuffer(nil)

	t := template.Must(template.New(n.AlarmId).Parse(tmpl))
	err := t.Execute(out, n)

	if err != nil {
		return "", &Err.Error{Code: Err.Econflict, Message: "Invalid template", Err: err}
	}
	return out.String(), nil
}

func recoverInvalidTemplate() {
	if err := recover(); err != nil {
		logrus.Error(fmt.Sprintf("Recovered from invalid template form: %s", err))
	}
}
