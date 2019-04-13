package notifications

import (
	"fmt"
	"testing"
	"time"
)

func TestNotification_ParseOk(t *testing.T) {
	tmpl := "Alarm {{.AlarmName}} has fired at {{.Timestamp}}"
	ts := time.Now()
	n := Notification{
		AlarmName: "test_alarm",
		Timestamp: ts,
	}

	data, err := n.Parse(tmpl)
	if err != nil {
		t.Error(err)
	}

	if data != fmt.Sprintf("Alarm %s has fired at %s", n.AlarmName, n.Timestamp.String()) {
		t.Error("Parsed text doesn't match")
	}
}

func TestNotification_ParseInvalidTemplate(t *testing.T) {
	tmpl := "Alarm {{.AlarmName}} has fired at {{.Timestamp}"
	ts := time.Now()
	n := Notification{
		AlarmName: "test_alarm",
		Timestamp: ts,
	}

	data, err := n.Parse(tmpl)
	if data != "" {
		t.Error("Got not empty result from invalid template.")
	}

	if err == nil {
		t.Error("Expected error from template parsing, but didn't get any.")
	}

}
