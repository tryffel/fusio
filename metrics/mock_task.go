package metrics

// Mock implementation of metrics task. Only implements Metrics interface
type MockTask struct {
}

func (m *MockTask) GaugeIncrease(name string, val float64) {
}

func (m *MockTask) GaugeDecrease(name string, val float64) {
}

func (m *MockTask) CounterIncrease(name string, val float64) {
}
