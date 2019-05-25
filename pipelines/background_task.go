package pipelines

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage"
	"sync"
	"time"
)

const (
	MinTimeInterval = 10 * time.Microsecond
	InputBufSize    = 20
)

// Delegate can import activities
type Delegate interface {
	// AddActivity adds value from either groups and optional device
	AddActivity(device string, groups *[]string, name string, value float64)
}

type Action struct {
	Device string
	Groups []string
	Name   string
	Value  float64
}

type BackgroundTask struct {
	lock             sync.RWMutex
	initialized      bool
	running          bool
	minInterval      time.Duration
	inMemoryInterval time.Duration
	inCacheInterval  time.Duration
	store            *storage.Store
	ticket           *time.Ticker

	runners       map[string]pipeline
	sigStop       chan bool
	inputActivity chan Action
}

func NewBackgroundTask(config *config.Config, store *storage.Store) (*BackgroundTask, error) {
	pt := &BackgroundTask{}

	pt.minInterval = time.Duration(config.Pipelines.MinInterval)
	pt.inMemoryInterval = time.Duration(config.Pipelines.InMemoryMaxInterval)
	pt.inCacheInterval = time.Duration(config.Pipelines.InCacheMaxInterval)
	pt.store = store

	if pt.minInterval < MinTimeInterval || pt.inMemoryInterval < MinTimeInterval || pt.inCacheInterval < MinTimeInterval {
		return pt, &Err.Error{Code: Err.Einternal,
			Err: errors.New(fmt.Sprintf("all intervals must be > %s", MinTimeInterval.String()))}
	}

	pt.sigStop = make(chan bool)
	pt.inputActivity = make(chan Action, InputBufSize)
	pt.initialized = true
	return pt, nil
}

func (b *BackgroundTask) Start() error {
	if !b.initialized {
		return errors.New("pipelines task not initialized correctly")
	}
	if b.running {
		return errors.New("pipelines task already running")
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	logrus.Debug("Starting pipelines task")
	b.running = true

	for _, v := range b.runners {
		go v.Run()
	}

	go b.loop()
	return nil
}

func (b *BackgroundTask) Stop() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.running = false
	b.sigStop <- true
}

func (bt *BackgroundTask) AddActivity(device string, groups *[]string, name string, value float64) {
	bt.inputActivity <- Action{Device: device, Groups: *groups, Name: name, Value: value}
}

func (bt *BackgroundTask) loop() {
	logrus.Info("Started pipelines task")
	for {
		select {
		case <-bt.sigStop:
			return
		case action := <-bt.inputActivity:
			logrus.Info(action.Name)

		}
	}
}
