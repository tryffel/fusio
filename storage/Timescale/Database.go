package Timescale

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	Err "github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/util"
	"time"
)

type Retention struct {
	name          string
	Duration      util.Interval
	Interval      util.Interval
	ChunkInterval util.Interval
}

// Validate config.Retentionpolicy into retention
func getRetention(p *config.RetentionPolicy) (Retention, error) {
	r := Retention{
		name:          "",
		Duration:      p.Duration,
		Interval:      p.Interval,
		ChunkInterval: p.ChunkInterval,
	}

	if r.Duration < r.ChunkInterval {
		return r, &Err.Error{Code: Err.Einvalid,
			Err: errors.New("retention chunk interval cannot be greater then retention duration")}
	}

	if time.Duration(p.Duration) < time.Hour {
		return r, &Err.Error{Code: Err.Einvalid, Err: errors.New("retention time cannot be < 1h")}
	}
	if time.Duration(p.Interval) > time.Duration(p.ChunkInterval) {
		return r, &Err.Error{Code: Err.Einvalid, Err: errors.New("retention bucket interval has to be longer " +
			"than retention interval")}
	}

	if time.Duration(p.Duration) < time.Hour*72 {
		r.name = fmt.Sprintf("ts_%dh", p.Duration.ToSeconds()/3600)
	} else {
		r.name = fmt.Sprintf("ts_%dd", p.Duration.ToSeconds()/3600/24)
	}
	return r, nil
}

type Engine struct {
	db         *gorm.DB
	retentions []Retention
}

func NewEngine(db *gorm.DB, config *config.TimeSeries) (*Engine, error) {
	engine := &Engine{db: db}

	ext, err := engine.ExtensionExists()
	if err != nil {
		return engine, Err.Wrap(&err, "could not initialize timescale engine")
	}
	if !ext {
		err = engine.EnableExtension()
		if err != nil {
			return engine, Err.Wrap(&err, "failed to enable timescale extension")
		}
	}
	if !ext {
		e := errors.New("TimescaleDB extension not enabled")
		return engine, Err.Wrap(&e, "")
	}

	for _, v := range config.Retentions {
		r, err := getRetention(&v)
		if err != nil {
			return engine, err
		}
		engine.retentions = append(engine.retentions, r)
	}

	return engine, nil
}

type Row struct {
	Value string
}

// Check timescale-extension exists in postgresql database
func (e *Engine) ExtensionExists() (bool, error) {
	var count uint8
	result := e.db.Raw("SELECT COUNT(extname) AS value FROM PG_EXTENSION WHERE extname='timescaledb'").Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	if count == 1 {
		return true, result.Error
	}
	return false, result.Error
}

// Enable timescale extension in database
func (e *Engine) EnableExtension() error {
	logrus.Warning("Enabling TimescaleDB extension")
	result := e.db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE")
	if result.Error == nil {
		logrus.Warn("TimescaleDB extension ok")
		return nil
	} else {
		if e, ok := result.Error.(*pq.Error); ok {
			switch e.Code {
			case "42501":
				logrus.Error("failed to create extension TimescaleDB due to insufficient permissions. " +
					"Superuser permissions are needed to create extension. You can also manually create extension by " +
					"running following query in database: 'CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;'")
				return &Err.Error{Code: Err.Econflict, Message: "Internal error", Err: e}
			}
		}
	}
	return nil
}
