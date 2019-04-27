package Timescale

import (
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/util"
	"testing"
	"time"
)

func TestGetRetentionValid(t *testing.T) {

	c := config.RetentionPolicy{
		Duration:      util.Interval(time.Hour * 24),
		Interval:      util.Interval(time.Second),
		ChunkInterval: util.Interval(time.Hour * 24),
	}

	r, err := getRetention(&c)
	if err != nil {
		t.Error("retention policy validation failed: ", err.Error())
	}

	if (&r).name != "ts_24h" {
		t.Errorf("invalid retention policy name, exptected 'ts_24h', got '%s'", r.name)
	}

	if r.Duration != c.Duration || r.Interval != c.Interval || r.ChunkInterval != c.ChunkInterval {
		t.Error("retention is incorrect")
	}

	c.Duration = util.Interval(time.Hour * 24 * 7)
	r, err = getRetention(&c)
	if err != nil {
		t.Error("retention policy validation failed: ", err.Error())
	}

	if (&r).name != "ts_7d" {
		t.Errorf("invalid retention policy name, exptected 'ts_7d', got '%s'", (&r).name)
	}
}

func TestGetRetentionInvalid(t *testing.T) {

	// Duration < ChunkInterval
	c := config.RetentionPolicy{
		Duration:      util.Interval(time.Hour * 24),
		Interval:      util.Interval(time.Second),
		ChunkInterval: util.Interval(time.Hour * 24 * 2),
	}

	_, err := getRetention(&c)
	if err == nil {
		t.Error("retention policy validation failed: invalid retention was accepted")
	}

	// duration < hour
	c.ChunkInterval = util.Interval(time.Hour * 24)
	c.Duration = util.Interval(time.Minute * 10)

	_, err = getRetention(&c)
	if err == nil {
		t.Error("retention policy validation failed: invalid retention was accepted")
	}

	// interval > bucketinterval
	c.ChunkInterval = util.Interval(time.Hour * 24)
	c.Interval = util.Interval(time.Hour * 24 * 3)

	_, err = getRetention(&c)
	if err == nil {
		t.Error("retention policy validation failed: invalid retention was accepted")
	}
}
