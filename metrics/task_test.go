package metrics

import (
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/storage"
	"github.com/tryffel/fusio/storage/repository_mock"
	"github.com/tryffel/fusio/util"
	"strings"
	"testing"
	"time"
)

func TestMetrics(T *testing.T) {

	// Setup

	logrus.SetLevel(logrus.ErrorLevel)

	store, _ := storage.NewMockStore()
	mockMeasurement := store.Measurement.(*repository_mock.MockMeasurementRepository)
	mockMeasurement.StoreResults = true

	conf := config.Config{}
	conf.Metrics.RunMetrics = true
	conf.Metrics.Interval = util.Interval(time.Millisecond * 1000)

	task, err := NewBackgroundTask(&conf, store)
	if err != nil {
		T.Error(err)
	}

	// Test
	// Push both counters and gauges into metrics task and check their values

	err = task.Start()
	if err != nil {
		T.Error(err)
		task.Stop()
		return
	}

	task.CounterIncrease("test_c_a", 1)
	task.CounterIncrease("test_c_b", 2)
	task.GaugeIncrease("test_g_a", 1)
	task.GaugeIncrease("test_g_b", 5)
	task.GaugeDecrease("test_g_b", 2)

	time.Sleep(time.Millisecond * 1100)

	results := mockMeasurement.Metrics

	if results["test_c_a_counter"] != float64(1) {
		T.Errorf("Metrics counter not working properly, exptected '%d', got '%f'", 1, results["test_c_a"])
	}
	if results["test_c_b_counter"] != float64(2) {
		T.Errorf("Metrics counter not working properly, exptected '%d', got '%f'", 2, results["test_c_b"])
	}
	if results["test_g_a_gauge"] != float64(1) {
		T.Errorf("Metrics gauge not working properly, exptected '%d', got '%f'", 1, results["test_g_a"])
	}
	if results["test_g_b_gauge"] != float64(3) {
		T.Errorf("Metri:wcs gauge not working properly, exptected '%d', got '%f'", 3, results["test_g_b"])
	}

	task.Stop()
}

func BenchmarkBackgroundTask_CounterIncrease_SingleMetrics(b *testing.B) {
	logrus.SetLevel(logrus.ErrorLevel)

	store, _ := storage.NewMockStore()
	mockMeasurement := store.Measurement.(*repository_mock.MockMeasurementRepository)
	mockMeasurement.StoreResults = false

	conf := config.Config{}
	conf.Metrics.RunMetrics = true
	conf.Metrics.Interval = util.Interval(time.Second * 600)

	task, err := NewBackgroundTask(&conf, store)
	if err != nil {
		b.Error(err)
	}

	err = task.Start()
	if err != nil {
		b.Error(err)
		task.Stop()
		return
	}

	for i := 0; i < b.N; i++ {
		task.CounterIncrease("benchmark_a", 1)
	}
}

func BenchmarkBackgroundTask_CounterIncrease_ThousandMetrics(b *testing.B) {
	logrus.SetLevel(logrus.ErrorLevel)

	store, _ := storage.NewMockStore()
	mockMeasurement := store.Measurement.(*repository_mock.MockMeasurementRepository)
	mockMeasurement.StoreResults = false

	conf := config.Config{}
	conf.Metrics.RunMetrics = true
	conf.Metrics.Interval = util.Interval(time.Second * 600)

	task, err := NewBackgroundTask(&conf, store)
	if err != nil {
		b.Error(err)
	}

	err = task.Start()
	if err != nil {
		b.Error(err)
		task.Stop()
		return
	}

	metricsNum := 1000

	for i := 0; i < b.N; i++ {
		s := strings.Builder{}
		s.WriteString(string(i % metricsNum))
		task.CounterIncrease("s", 1)
	}
}
