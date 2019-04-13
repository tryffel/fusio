package alarm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/metrics"
	"github.com/tryffel/fusio/storage"
	"runtime/debug"
	"sync"
	"time"
)

const (
	MinTimeInterval = time.Second * 5
)

// Background task to run periodically and valuate alarms
type BackgroundTask struct {
	lock        sync.RWMutex
	initialized bool
	running     bool
	interval    time.Duration
	lastRun     time.Time
	store       *storage.Store
	metrics     metrics.Metrics
}

// NewBackgrounTask Create new alarming background task
func NewBackgroundTask(config config.Config, store *storage.Store, metric metrics.Metrics) (*BackgroundTask, error) {
	bt := &BackgroundTask{}
	interval := time.Duration(config.Alarms.Interval)
	if interval < MinTimeInterval {
		return bt, errors.New(fmt.Sprintf("Minimum interval for background tasks is %s", MinTimeInterval.String()))
	}
	bt.interval = interval
	bt.store = store
	bt.metrics = metric
	bt.initialized = true
	return bt, nil
}

// Start task
func (b *BackgroundTask) Start() error {
	if !b.initialized {
		return &Err.Error{Code: Err.Einternal, Err: errors.New("alarm task not initialized properly")}
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.running {
		logrus.Debug("Starting alarm task")
		b.running = true
		go b.loop()
	}
	return nil
}

// Stop task
func (b *BackgroundTask) Stop() {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.running {
		logrus.Debug("Stopping alarm task")
		b.running = false
	}
}

// Alarming loop
func (b *BackgroundTask) loop() {
	logrus.Info("Running alarms every ", b.interval.String())
	for b.IsRunning() {
		b.runAlarms()
		time.Sleep(b.interval)
	}
	logrus.Info("Alarm task stopped")
}

// IsRunning check if task is running
func (b *BackgroundTask) IsRunning() bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.running
}

// Run alarms evaluation
func (b *BackgroundTask) runAlarms() {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(*Err.Error); ok {
				e.Wrap("Panic in alarm task")
				logrus.Error(e.Error())
			} else {
				logrus.Error("Panic in alarm background task: ", err)
			}
			debug.PrintStack()
			time.Sleep(time.Second * 15)
		}
	}()

	alarms, err := b.store.Alarm.GetAlarmsToValuate(b.interval)
	if err != nil {
		if e, ok := err.(*Err.Error); ok {
			e.Wrap("Failed to gather alarms to valuate")
			logrus.Error(e)
		}
		logrus.Error(errors.Wrap(err, "Failed to gather alarms to valuate"))
		return
	}
	logrus.Debug("Evaluating ", len(*alarms), " alarms")
	b.metrics.CounterIncrease("alarm_evaluate", float64(len(*alarms)))
	if len(*alarms) == 0 {
		return
	} else {
		RunAlarms(*alarms, *b.store, b.metrics)
	}
	duration := time.Since(start)
	b.metrics.CounterIncrease("alarm_evaluation_time_us", float64(duration.Nanoseconds()/1000))
}
