package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Interval time.duration that marshals to seconds and is able to unmarshal many formats
type Interval time.Duration

func (i Interval) MarshalYAML() (interface{}, error) {
	return time.Duration(i).String(), nil
}

func (i *Interval) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var input string
	if err := unmarshal(&input); err != nil {
		return err
	}

	duration, err := time.ParseDuration(input)
	if err != nil {
		return err
	}

	*i = Interval(duration)
	return nil
}

func (i *Interval) ToDuration() time.Duration {
	return time.Duration(*i)
}

func NewInterval(d time.Duration) Interval {
	return Interval(d)
}

func (i *Interval) ToSeconds() int64 {
	return int64(i.ToDuration().Seconds())
}

func (i *Interval) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.0f", i.ToDuration().Seconds())), nil
}

func (i *Interval) UnmarshalJSON(b []byte) error {
	d := time.Duration(10)

	str := string(b)
	if _, err := strconv.Atoi(str); err == nil {
		err = json.Unmarshal(b, &d)
	} else {
		d, err = time.ParseDuration(strings.Trim(string(str), "\""))
		if err != nil {
			return err
		}
	}

	*i = Interval(d)
	return nil
}
