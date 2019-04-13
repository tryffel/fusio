// Influxdb driver
// All measurements are saved as follows
// field value is always 'measurementValue' to ease continuous queries
// tags are device=id, groups=groupB<separator>groupB<separator>..., 'measurementKey'=measurement_name
//

package Influxdb

import (
	"encoding/json"
	"errors"
	"fmt"
	influx_client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

const (
	measurementName  = "measurement"
	measurementValue = "value"
	measurementKey   = "key"
	deviceName       = "device"
	groupName        = "group"
	groupSeparator   = ";"
	// Maximum size to return measurements for
	maxHistorySize = 300

	metricsMeasurement = "metrics"
	metricsName        = "name"
)

// Public interface for influxdb client
type Client interface {
	// Write measurements for given device and groups
	// Measurements is map of measurement name to measurement value.
	Write(device string, groups []string, m Measurements) error
	// Read measurement return array of measurements as defined in inputs array. Measurements are
	// gathered between from and to timestamps and max length is of n
	Read(device string, group string, filters []Filter, from time.Time, to time.Time, n int64) (Batch, error)

	// Write metrics for given name
	WriteMetrics(name string, value float64) error

	WriteMetricsBatch(batch *map[string]float64) error

	// Get Measurements for device
	GetDeviceMeasurements(device string) ([]string, error)

	// Get measurements for group
	GetGroupMeasurements(group string) ([]string, error)
}

type client struct {
	client      influx_client.Client
	db          string
	retentions  []retention
	sendMetrics bool
	logQueries  bool
	logger      *util.SqlLogger
	databaseManager
}

func (c *client) GetDeviceMeasurements(device string) ([]string, error) {
	query := fmt.Sprintf(`SHOW TAG VALUES WITH key="%s" WHERE "%s"='%s'`, measurementKey, deviceName, device)
	q := influx_client.NewQuery(query, c.db, "")

	res, err := c.client.Query(q)
	if err != nil {
		c.logQuery(query, err, nil)
		return []string{}, err
	}
	if res.Error() != nil {
		c.logQuery(query, res.Error(), &res.Results[0])
	}

	return keyValueToArray(&res.Results)
}

func (c *client) GetGroupMeasurements(group string) ([]string, error) {
	query := fmt.Sprintf(`SHOW TAG VALUES WITH key="%s" WHERE "%s"=~/.*%s.*/`, measurementKey, groupName, group)
	q := influx_client.NewQuery(query, c.db, "")

	res, err := c.client.Query(q)
	if err != nil {
		c.logQuery(query, err, nil)
		return []string{}, err
	}
	if res.Error() != nil {
		c.logQuery(query, res.Error(), &res.Results[0])
	}

	if len(res.Results) == 0 {
		return []string{}, nil
	}

	return keyValueToArray(&res.Results)
}

func (c *client) Write(device string, groups []string, measurements Measurements) error {
	batch, err := influx_client.NewBatchPoints(influx_client.BatchPointsConfig{Database: c.db})
	if err != nil {
		return err
	}

	for name, v := range measurements {
		tags := map[string]string{
			groupName:      strings.Join(groups, groupSeparator),
			deviceName:     device,
			measurementKey: name,
		}

		fields := map[string]interface{}{
			measurementValue: v.Value,
		}
		point, err := influx_client.NewPoint(measurementName, tags, fields, v.Timestamp)
		if err != nil {
			return err
		}
		batch.AddPoint(point)
	}
	err = c.client.Write(batch)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Read(device string, group string, filters []Filter, from time.Time, to time.Time, n int64) (Batch, error) {

	if from.Nanosecond() > to.Nanosecond() {
		return Batch{}, errors.New("influxdb query time range has to be positive")
	}

	// Proper retention to match time interval
	retention, err := c.getRetentionPolicy(from, to)
	if err != nil {
		return Batch{}, err
	}

	// Construct separate query for each input based on base_query. Query them as batch and parse result into Batch
	baseQuery := `SELECT %s AS %s FROM "%s"."%s" WHERE %s AND %s GROUP BY %s fill(none) limit %d`
	fullQuery := ""
	timeQuery := fmt.Sprintf("time <= %ds AND time >= %ds", to.Unix(), from.Unix())
	groupQuery := fmt.Sprintf("time(%ds)", int64(getGroupByTime(from, to, n, *retention).Seconds()))
	limit := n
	if n > maxHistorySize {
		limit = maxHistorySize
	}

	deviceGroup := getDeviceGroupClause(device, group)
	whereClause := deviceGroup
	if whereClause != "" {
		whereClause = fmt.Sprintf(" %s AND ", whereClause)
	}
	whereClause = fmt.Sprintf("%s %s", whereClause, timeQuery)

	for _, filter := range filters {
		measurementQuery := fmt.Sprintf(`"%s"='%s'`, measurementKey, filter.Key)
		query := fmt.Sprintf(baseQuery, filter.influxString(), filter.StringSimplified(), retention.name, measurementName, whereClause, measurementQuery, groupQuery, limit)
		if fullQuery != "" {
			fullQuery = fmt.Sprintf("%s; %s", fullQuery, query)
		} else {
			fullQuery = query
		}
	}

	q := influx_client.NewQuery(fullQuery, c.db, "s")

	res, err := c.client.Query(q)

	if res == nil {
		logrus.Error("Influxdb query returned empty response")
		return Batch{}, errors.New("empty result from influxdb")
	}

	if err != nil {
		c.logQuery(fullQuery, err, &res.Results[0])
	}
	c.logQuery(fullQuery, res.Error(), &res.Results[0])

	measurements, err := seriesToMeasurements(&res.Results, filters)
	if err != nil {
		return Batch{}, nil
	}
	logrus.Debug(measurements)

	return *measurements, nil

}

// WriteMetrics writes metrics if enabled in config
func (c *client) WriteMetrics(name string, value float64) error {
	if !c.sendMetrics {
		return nil
	}

	tags := map[string]string{
		metricsName: name,
	}

	fields := map[string]interface{}{
		"value": value,
	}

	point, err := influx_client.NewPoint(metricsMeasurement, tags, fields, time.Now())
	if err != nil {
		return err
	}

	batch, err := influx_client.NewBatchPoints(influx_client.BatchPointsConfig{Database: c.db})
	if err != nil {
		return err
	}

	batch.AddPoint(point)
	return c.client.Write(batch)
}

func (c *client) WriteMetricsBatch(batch *map[string]float64) error {
	if !c.sendMetrics {
		return nil
	}

	if len(*batch) == 0 {
		return nil
	}

	batchPoint, err := influx_client.NewBatchPoints(influx_client.BatchPointsConfig{Database: c.db})
	if err != nil {
		return err
	}

	for i, v := range *batch {
		tags := map[string]string{
			metricsName: i,
		}
		fields := map[string]interface{}{
			measurementValue: v,
		}

		point, err := influx_client.NewPoint(metricsMeasurement, tags, fields, time.Now())
		if err != nil {
			logrus.Error(err)
		}

		batchPoint.AddPoint(point)
	}
	return c.client.Write(batchPoint)
}

func NewClient(config *config.Influxdb, sendMetrics bool, logQueries bool, logger *util.SqlLogger) (Client, error) {
	c := &client{}

	influx, err := influx_client.NewHTTPClient(influx_client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", config.Host, config.Port)})
	if err != nil {
		return c, err
	}

	c.client = influx
	c.db = config.Database
	c.retentions = retentions
	c.sendMetrics = sendMetrics
	c.logQueries = logQueries
	c.logger = logger

	duration, version, err := c.client.Ping(time.Second * 10)
	if err != nil {
		logrus.Error("Failed to connect to influxdb: ", err)
		return c, err
	}

	logrus.Debugf("Connected to influxdb, ping: %d ms, version: %s", duration.Nanoseconds()/1000000, version)

	err = InitDB(c)
	if err != nil {
		logrus.Error(err)
		return c, err
	}

	return c, nil
}

// logQuery logs query q, error and number of rows returner
func (c *client) logQuery(q string, err error, series *influx_client.Result) {
	if err != nil {
		logrus.Errorf("Influxdb error: %s query: %s ", err, q)
	}

	if c.logQueries {
		seriesNuM := len(series.Series)
		var rowNum = 0
		if seriesNuM > 0 {
			rowNum = len(series.Series[0].Values)
		}

		fields := map[string]interface{}{
			"query":  q,
			"series": seriesNuM,
			"rows":   rowNum,
		}
		c.logger.WithFields(fields).Info("influxdb query")
	}
}

// Construct device and group filters for influxdb query. If both are empty, empty query is returned
func getDeviceGroupClause(deviceId string, groupId string) string {
	if deviceId == "" && groupId == "" {
		return ""
	}

	query := ""
	if deviceId != "" {
		query = fmt.Sprintf(`"%s"='%s'`, deviceName, deviceId)
		if groupId != "" {
			query = fmt.Sprintf("%s AND ", query)
		}
	}

	if groupId != "" {
		query = fmt.Sprintf(`%s "%s"=~/.*%s.*/`, query, groupName, groupId)
	}
	return query
}

// getGroupByTime constructs time interval for given times and retention policy
func getGroupByTime(from time.Time, to time.Time, n int64, retention retention) time.Duration {
	duration := to.Sub(from)

	if n > maxHistorySize {
		n = maxHistorySize
	}
	interval := time.Duration(duration.Nanoseconds() / n)

	if interval < retention.samplingRate {
		logrus.Debugf("Querying with samplerate of %s, "+
			"but downsampling rate is %s, using downsampled rate", interval.String(), retention.samplingRate.String())
		return retention.samplingRate
	}
	return interval
}

// Results are expected to be on same order as filters
func seriesToMeasurements(measurements *[]influx_client.Result, filters []Filter) (*Batch, error) {
	if len(*measurements) != len(filters) {
		return &Batch{}, errors.New("mismatch of inputs and results")
	}
	result := &Batch{}

	keys := make([]string, len(filters))

	for i, v := range filters {
		keys[i] = v.StringSimplified()
	}

	for i, measurement := range *measurements {
		if len(measurement.Series) > 0 {
			columns := columnsAsMap(measurement.Series[0])
			name := filters[i].StringSimplified()

			series := make([]Point, len(measurement.Series[0].Values))

			for index, point := range measurement.Series[0].Values {
				num, _ := point[columns["time"]].(json.Number).Int64()
				var value float32 = 0.0
				if point[columns[name]] != nil {
					val, _ := point[columns[name]].(json.Number).Float64()
					value = float32(val)
				}

				p := Point{
					Timestamp: time.Unix(num, 0),
					Value:     value,
				}
				series[index] = p
			}
			(*result)[name] = series
		}
	}

	return result, nil
}

func keyValueToArray(measurements *[]influx_client.Result) ([]string, error) {
	if len(*measurements) == 0 {
		return []string{}, nil
	}

	if (*measurements)[0].Series == nil {
		return []string{}, nil

	}

	series := (*measurements)[0].Series[0]
	array := make([]string, len(series.Values))
	for i, v := range series.Values {
		array[i] = v[1].(string)
	}

	return array, nil
}
