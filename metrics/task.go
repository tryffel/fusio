package metrics

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/storage"
	"strings"
	"sync"
	"time"
)

const (
	MinTimeInterval = time.Second
	chanBufferSize  = 20
)

// Metrics interface provides methods to store metrics
type Metrics interface {
	// Gauge can be both increased and decreased
	GaugeIncrease(name string, val float64)
	GaugeDecrease(name string, val float64)
	// Counter can be only increased
	CounterIncrease(name string, val float64)
}

type metricsAction struct {
	name string
	val  float64
}

type BackgroundTask struct {
	lock        sync.RWMutex
	initialized bool
	running     bool
	interval    time.Duration
	lastRun     time.Time
	store       *storage.Store
	ticker      *time.Ticker

	counters map[string]float64
	gauges   map[string]float64

	counterBuf chan metricsAction
	gaugeBuf   chan metricsAction
}

func (b *BackgroundTask) GaugeIncrease(name string, val float64) {
	b.gaugeBuf <- metricsAction{name: name, val: val}
}

func (b *BackgroundTask) GaugeDecrease(name string, val float64) {
	b.gaugeBuf <- metricsAction{name: name, val: -val}
}

func (b *BackgroundTask) CounterIncrease(name string, val float64) {
	b.counterBuf <- metricsAction{name: name, val: val}
}

func NewBackgroundTask(config *config.Config, store *storage.Store) (*BackgroundTask, error) {
	bt := &BackgroundTask{}
	bt.interval = time.Duration(config.Metrics.Interval)
	if bt.interval < MinTimeInterval {
		return bt, errors.New(fmt.Sprintf("minimum interval for metrics is %s", MinTimeInterval.String()))
	}
	bt.store = store
	bt.ticker = time.NewTicker(bt.interval)

	bt.counters = make(map[string]float64, 0)
	bt.gauges = make(map[string]float64)
	bt.counterBuf = make(chan metricsAction, chanBufferSize)
	bt.gaugeBuf = make(chan metricsAction, chanBufferSize)
	bt.initialized = true
	return bt, nil
}

func (b *BackgroundTask) ensureCounter(name string) {
	_, found := b.counters[name]

	if found {
		return
	}
	b.counters[name] = float64(0)
}

func (b *BackgroundTask) ensureGauge(name string) {
	_, found := b.gauges[name]
	if found {
		return
	}
	b.gauges[name] = float64(0)
}

func (b *BackgroundTask) Start() error {
	if !b.initialized {
		return errors.New("metrics task not initialized correctly")
	}
	if b.running {
		return errors.New("metrics task already running")
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	logrus.Debug("Starting metrics task")
	b.running = true
	go b.loop()
	return nil
}

func (b *BackgroundTask) Stop() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.running = false
}

func (b *BackgroundTask) loop() {
	logrus.Debug("Pushing metrics every ", b.interval.String())
	for b.running {
		select {
		case item := <-b.counterBuf:
			b.ensureCounter(item.name)
			b.counters[item.name] += item.val
		case item := <-b.gaugeBuf:
			b.ensureGauge(item.name)
			b.gauges[item.name] += item.val
		case <-b.ticker.C:

			measurements := make(map[string]float64)

			for i, v := range b.counters {
				var name strings.Builder
				name.WriteString(i)
				name.WriteString("_")
				name.WriteString("counter")
				measurements[name.String()] = v
			}

			for i, v := range b.gauges {
				var name strings.Builder
				name.WriteString(i)
				name.WriteString("_")
				name.WriteString("gauge")
				measurements[name.String()] = v
			}

			go b.pushMetrics(&measurements)
		}
	}
	logrus.Info("Stopping metrics task")
}

func (b *BackgroundTask) pushMetrics(metrics *map[string]float64) {

	if len(*metrics) == 0 {
		return
	}

	logrus.Debug("Pushing ", len(*metrics), " metrics")

	start := time.Now()
	err := b.store.Measurement.WriteMetricsBatch(metrics)
	if err != nil {
		logrus.Error(err)
	}

	took := time.Since(start)

	logrus.Debug("Pushing metrics took ", took.String())
}
