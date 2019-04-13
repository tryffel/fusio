package Influxdb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Match 'mean(measurement_test)'
var regexSelector, _ = regexp.Compile(`^([a-zA-Z_]+)\(([a-zA-Z_]+)\)$`)

// Matches moving_average(mean(usage_system),10)
var regexSelectorAndTransformation, _ = regexp.Compile(`^(\w+)?\((\w+)\((\w+)\),?(\d+)?\)$`)

var regexOperators, _ = regexp.Compile(`[-+=><*/]`)

// Information about the contents of filter.
const (
	filterSimple               = 0
	filterTransformed          = 1
	filterTransformedParameter = 2
)

type Filter struct {
	Selector       string `json:"selector"`
	Transformation string `json:"transformation"`
	Key            string `json:"key"`
	// Additional parameters for transformations
	Param      int64 `json:"param"`
	FilterType int   `json:"type"`
}

// String get filter as original string
func (f *Filter) String() string {
	if f.FilterType == filterSimple {
		return fmt.Sprintf("%s(%s)", f.Selector, f.Key)
	}
	if f.FilterType == filterTransformed {
		return fmt.Sprintf("%s(%s(%s))", f.Transformation, f.Selector, f.Key)
	}
	if f.FilterType == filterTransformedParameter {
		return fmt.Sprintf("%s(%s(%s),%d)", f.Transformation, f.Selector, f.Key, f.Param)
	}
	return "Unknown filter type"
}

// StringEscaped get filter as escaped string, for influx queries
func (f *Filter) StringEscaped() string {
	if f.FilterType == filterSimple {
		return fmt.Sprintf("%s(\"%s\")", f.Selector, f.Key)
	}
	if f.FilterType == filterTransformed {
		return fmt.Sprintf("%s(%s(\"%s\"))", f.Transformation, f.Selector, f.Key)
	}
	if f.FilterType == filterTransformedParameter {
		return fmt.Sprintf("%s(%s(\"%s\"),%d)", f.Transformation, f.Selector, f.Key, f.Param)
	}
	return "Unknown filter type"
}

// influxString gets filter formatted as 'mean("value")', where value is hardcoded 'measurementValue*
func (f *Filter) influxString() string {
	if f.FilterType == filterSimple {
		return fmt.Sprintf("%s(\"%s\")", f.Selector, measurementValue)
	}
	if f.FilterType == filterTransformed {
		return fmt.Sprintf("%s(%s(\"%s\"))", f.Transformation, f.Selector, measurementValue)
	}
	if f.FilterType == filterTransformedParameter {
		return fmt.Sprintf("%s(%s(\"%s\"),%s)", f.Transformation, f.Selector, f.Key, measurementValue)
	}
	return "Unknown filter type"
}

// StringSimplified get string as a placeholder: mean_max_temperature
func (f *Filter) StringSimplified() string {
	if f.FilterType == filterSimple {
		return fmt.Sprintf("%s_%s", f.Selector, f.Key)
	}
	if f.FilterType == filterTransformed || f.FilterType == filterTransformedParameter {
		return fmt.Sprintf("%s_%s_%s", f.Transformation, f.Selector, f.Key)
	}
	return "Unknown filter type"
}

// FilterFromString Attempt to validate filter in string and return parsed filters or error with message description
// Valid string is: mean(temperature) - derivative(mean(temperature),10) > 10
// Invalid string is: max mean(temperature) is 0
func FilterFromString(s string) (*[]Filter, error) {
	var err error
	s = strings.Replace(s, " ", "", -1)
	sub := regexOperators.Split(s, -1)

	filter := make([]Filter, 0)

	for _, v := range sub {
		match := regexSelectorAndTransformation.FindStringSubmatch(v)

		// Match simple filter 'mean(measurement)'
		if match != nil {
			f := Filter{
				Selector:       match[2],
				Key:            match[3],
				Transformation: match[1],
				Param:          -1,
				FilterType:     filterTransformed,
			}

			// Match possible integer parameter
			if len(match[4]) > 0 {
				f.Param, err = strconv.ParseInt(match[4], 10, 64)
				f.FilterType = filterTransformedParameter
				if err != nil {
					return &filter, errors.New(fmt.Sprintf("'%s' not integer", match[4]))
				}
			}
			filter = append(filter, f)
		} else {
			match = regexSelector.FindStringSubmatch(v)
			// Match filter 'derivative(mean(measurement),10)
			if match != nil {
				f := Filter{
					Selector:   match[1],
					Key:        match[2],
					FilterType: filterSimple,
				}
				filter = append(filter, f)
			}
		}
	}
	if len(filter) == 0 {
		return &filter, errors.New("Filters can't be empty")

	}
	return &filter, nil
}
