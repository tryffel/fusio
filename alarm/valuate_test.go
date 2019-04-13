package alarm

import (
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/models"
	"testing"
	"time"
)

func TestValuateAlarmSeries(T *testing.T) {

	var rounds = 5

	raw_expression := "mean(temperature) + max(temperature) > 35"
	simplified_expression := "mean_temperature + max_temperature > 35"
	filters, err := Influxdb.FilterFromString(raw_expression)
	if err != nil {
		T.Errorf("Failed to create influxdb filters from string: %s", err)
	}

	query := &models.AlarmQuery{
		Filters:    *filters,
		Interval:   time.Second * 1,
		Expression: simplified_expression,
		Limit:      int64(rounds),
	}

	batch := Influxdb.Batch{}

	tempsArr := make([]Influxdb.Point, rounds)
	for i := 0; i < rounds; i++ {
		tempsArr[i] = Influxdb.Point{
			Timestamp: time.Now().Add(-time.Second * 5),
			Value:     20.0,
		}
	}

	batch["mean_temperature"] = tempsArr
	batch["max_temperature"] = tempsArr

	// Should give positive result
	status, err, _ := ValuateSeries(query, batch)

	if !status {
		T.Errorf("Failed to fire alarm: %s", query.Expression)
	}

	if err != nil {
		T.Error(err)
	}

	// Should give negative result
	batch["max_temperature"][2].Value = 0
	status, err, _ = ValuateSeries(query, batch)

	if status {
		T.Errorf("False firing alarm: %s", query.Expression)
	}

	if err != nil {
		T.Error(err)
	}

}

func BenchmarkValuateTwoSeriesFiveRounds(b *testing.B) {

	rounds := 5

	raw_expression := "mean(temperature) + max(temperature) > 35"
	simplified_expression := "mean_temperature + max_temperature > 35"
	filters, err := Influxdb.FilterFromString(raw_expression)
	if err != nil {
		b.Errorf("Failed to create influxdb filters from string: %s", err)
	}

	query := &models.AlarmQuery{
		Filters:    *filters,
		Interval:   time.Second * 1,
		Expression: simplified_expression,
		Limit:      int64(rounds),
	}

	batch := Influxdb.Batch{}

	tempsArr := make([]Influxdb.Point, rounds)
	for i := 0; i < rounds; i++ {
		tempsArr[i] = Influxdb.Point{
			Timestamp: time.Now().Add(-time.Second * 5),
			Value:     20.0,
		}
	}

	batch["mean_temperature"] = tempsArr
	batch["max_temperature"] = tempsArr

	for i := 0; i < b.N; i++ {
		// Should give positive result
		ValuateSeries(query, batch)
	}
}

func BenchmarkValuateTwoSeriesSingleRound(b *testing.B) {

	rounds := 1

	raw_expression := "mean(temperature) + max(temperature) > 35"
	simplified_expression := "mean_temperature + max_temperature > 35"
	filters, err := Influxdb.FilterFromString(raw_expression)
	if err != nil {
		b.Errorf("Failed to create influxdb filters from string: %s", err)
	}

	query := &models.AlarmQuery{
		Filters:    *filters,
		Interval:   time.Second * 1,
		Expression: simplified_expression,
		Limit:      int64(rounds),
	}

	batch := Influxdb.Batch{}

	tempsArr := make([]Influxdb.Point, rounds)
	for i := 0; i < rounds; i++ {
		tempsArr[i] = Influxdb.Point{
			Timestamp: time.Now().Add(-time.Second * 5),
			Value:     20.0,
		}
	}

	batch["mean_temperature"] = tempsArr
	batch["max_temperature"] = tempsArr

	for i := 0; i < b.N; i++ {
		// Should give positive result
		ValuateSeries(query, batch)
	}
}
