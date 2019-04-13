package metrics

import (
	"time"
)

type Timer struct {
	started time.Time
	stopped time.Time
	Name    string
}

func NewTimer(name string) Timer {
	return Timer{Name: name}
}

func (t *Timer) Start() {
	t.started = time.Now()
}

func (t *Timer) Stop() {
	t.stopped = time.Now()
}

func (t *Timer) Duration() time.Duration {
	return t.stopped.Sub(t.started)
}
