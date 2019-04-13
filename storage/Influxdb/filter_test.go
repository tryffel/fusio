package Influxdb

import (
	"testing"
)

func TestParseSingleInfluxFilter(t *testing.T) {
	validFilters := []struct {
		in          string
		original    string
		escaped     string
		placeholder string
	}{
		{"mean(temperature) > 10", "mean(temperature)", "mean(\"temperature\")", "mean_temperature"},
		{"max(input) = -5", "max(input)", "max(\"input\")", "max_input"},
		{"diff(mean(temp),10) > 10", "diff(mean(temp),10)", "diff(mean(\"temp\"),10)", "diff_mean_temp"},
	}

	invalidFilters := []string{
		" a - diff(aabc > 0",
		"b=2 == max a",
	}

	for _, c := range validFilters {
		got, err := FilterFromString(c.in)
		if err != nil {
			t.Error(err)
		}
		if len(*got) != 1 {
			t.Error("Received too many or little filters!")
		}
		if (*got)[0].String() != c.original {
			t.Errorf("Original string doesn't match filters string output: %s, got %s", c.original, (*got)[0].String())
		}
		if (*got)[0].StringEscaped() != c.escaped {
			t.Errorf("Escaped string doesn't match, %s, got %s", c.original, (*got)[0].StringEscaped())
		}
		if (*got)[0].StringSimplified() != c.placeholder {
			t.Errorf("Simplified string doesn't match, %s, got %s", c.original, (*got)[0].StringSimplified())
		}
	}

	for _, c := range invalidFilters {
		got, err := FilterFromString(c)
		if err == nil {
			t.Errorf("No error raised on invalid filter: %v", got)
		}
	}
}

func TestParseMultipleInfluxFilters(t *testing.T) {
	input := "mean(temp) - max(temp) > diff(mean(max),10)"
	strings := [3]string{}
	strings[0] = "mean(temp)"
	strings[1] = "max(temp)"
	strings[2] = "diff(mean(max),10)"

	escaped := [3]string{}
	escaped[0] = "mean(\"temp\")"
	escaped[1] = "max(\"temp\")"
	escaped[2] = "diff(mean(\"max\"),10)"

	simplified := [3]string{}
	simplified[0] = "mean_temp"
	simplified[1] = "max_temp"
	simplified[2] = "diff_mean_max"

	filters, err := FilterFromString(input)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for i, f := range *filters {
		if f.String() != strings[i] {
			t.Errorf("Original string doesn't match: %s, got: %s", strings[i], f.String())
		}
		if f.StringEscaped() != escaped[i] {
			t.Errorf("Escaped string doesn't match: %s, got: %s", escaped[i], f.StringEscaped())
		}
		if f.StringSimplified() != simplified[i] {
			t.Errorf("Simplified string doens't match: %s, got %s", simplified[i], f.StringSimplified())
		}
	}
}

func BenchmarkSingleInfluxFilter(b *testing.B) {
	filter := "mean(temperature) > 10 "

	for n := 0; n < b.N; n++ {
		FilterFromString(filter)
	}
}

func BenchmarkFilter_StringSimplified(b *testing.B) {
	filter := "mean(temperature) > 10"
	f, err := FilterFromString(filter)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		(*f)[0].StringSimplified()
	}
}

func BenchmarkFilter_String(b *testing.B) {
	filter := "mean(temperature) > 10"
	f, err := FilterFromString(filter)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		(*f)[0].String()
	}
}

func BenchmarkFilter_StringEscaped(b *testing.B) {
	filter := "mean(temperature) > 10"
	f, err := FilterFromString(filter)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		(*f)[0].StringEscaped()
	}
}
