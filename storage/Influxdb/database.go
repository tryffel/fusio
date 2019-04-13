package Influxdb

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	influx_client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"time"
)

// Retention represents Retention policy in influxdb
type retention struct {
	// Primary: default retention when inserting new measurements.
	// Only one retention can be default
	primary bool
	name    string
	// Duration how long data will retain in policy
	duration time.Duration
	// samplingRate Rate when downsampling measurements
	samplingRate time.Duration
}

// Default retention policies
// Retentions MUST be in ascending order by duration
var retentions = []retention{
	// 1 day of original measurements
	{
		primary:  true,
		name:     "1-day",
		duration: time.Hour * 24,
	},
	{
		primary:      false,
		name:         "6-months",
		duration:     time.Hour * 24 * 182,
		samplingRate: time.Minute * 30,
	},
	{
		primary:      false,
		name:         "3-years",
		duration:     time.Hour * 24 * 365 * 3,
		samplingRate: time.Hour * 24,
	},
}

const (
	ContinuousQueryPrefix = "cq_to_"
)

// DatabaseManager interface that can manage database
type databaseManager interface {
	CreateRetentionPolicies() error
	CreateContinuousQueries() error
	RetentionPoliciesExist() (bool, error)
	ContinuousQueriesExist() (bool, error)
	DatabaseExists() (bool, error)
	CreateDatabase() error
}

// createRetentionPolicies creates all RPs defined in client.
// This doesn't check for existing RPs, but returns error in case of failure
// function cannot handle if only one RP from many is missing, as it presumes either none or all RPs are present
func (c *client) CreateRetentionPolicies() error {
	if len(c.retentions) == 0 {
		return errors.New("no retention policies defined")
	}
	primaryExists := false

	logrus.Info("Creating influxdb retention policies ")
	for _, v := range c.retentions {
		query := fmt.Sprintf("CREATE RETENTION POLICY \"%s\" ON \"%s\" DURATION %ds REPLICATION 1",
			v.name, c.db, int64(v.duration.Seconds()))
		if v.primary {
			if primaryExists {
				err := errors.New("only one default retention police is allowed")
				return err
			}
			primaryExists = true
			query = fmt.Sprintf("%s DEFAULT", query)
			logrus.Debug(query)
		}
		q := influx_client.NewQuery(query, "", "")
		res, err := c.client.Query(q)
		if err != nil {
			return err
		}
		if res.Error() != nil {
			return res.Error()
		}

	}
	return nil
}

func (c *client) RetentionPoliciesExist() (bool, error) {
	query := influx_client.NewQuery("SHOW RETENTION POLICIES", c.db, "")
	res, err := c.client.Query(query)
	if err != nil {
		return false, err
	}
	if res.Error() != nil {
		return false, res.Error()
	}

	result := res.Results[0].Series[0]
	columnsMap := make(map[string]int)

	// Get column names
	for i, v := range result.Columns {
		columnsMap[v] = i
	}

	// Get configured retention policies as array
	retentionsMap := make(map[string]bool, len(c.retentions))
	for _, v := range c.retentions {
		retentionsMap[v.name] = true
	}

	// Check existing retentions against configured retentions and log if there are non-configured RPs
	// If found RPs name matches configured, delete it from retentionsMap
	for _, v := range result.Values {

		if retentionsMap[v[columnsMap["name"]].(string)] == false {
			logrus.Warnf("Found unused retention policy in influxdb: %s", v[columnsMap["name"]].(string))
		}
		delete(retentionsMap, v[columnsMap["name"]].(string))
	}

	// If map is still not empty, some RP is missing

	if len(retentionsMap) == len(c.retentions) {
		return false, nil
	}

	if len(retentionsMap) > 0 {
		return false, errors.New("some retention policies are missing in influxdb")
	}

	return true, nil

}

func (c *client) CreateContinuousQueries() error {
	format := `CREATE CONTINUOUS QUERY "cq_to_%s" on "%s"
BEGIN SELECT mean(%s) AS %s  
INTO "%s"."%s".:
MEASUREMENT FROM "%s"./.*/ 
GROUP BY time(%ds), * END`

	for i, v := range c.retentions {
		// Create only on non-default RP
		if !v.primary && i >= 1 {
			query := fmt.Sprintf(format, v.name, c.db, measurementValue, measurementValue, c.db, v.name, c.retentions[i-1].name, int64(v.samplingRate.Seconds()))
			logrus.Debug(query)
			res, err := c.client.Query(influx_client.NewQuery(query, "", ""))
			if err != nil || res.Error() != nil {
				logrus.Error("Failed to create continuous queries on influx: ", err, res.Error())
				return err
			}
		}
	}
	return nil
}

func (c *client) ContinuousQueriesExist() (bool, error) {
	// No need for CQs if only one RP
	if len(c.retentions) <= 1 {
		return true, nil
	}

	query := influx_client.NewQuery("SHOW CONTINUOUS QUERIES", c.db, "")
	res, err := c.client.Query(query)
	if err != nil {
		return false, err
	}
	if res.Error() != nil {
		return false, res.Error()
	}

	// Get correct series
	index := -1
	for i, v := range res.Results[0].Series {
		if v.Name == c.db {
			index = i
		}
	}

	if index < 0 {
		return false, errors.New("could not find influxdb continuous queries for given database")
	}

	columnsMap := columnsAsMap(res.Results[0].Series[index])

	// Map for all CQs wanted
	queriesMap := make(map[string]bool, len(c.retentions)-1)
	for _, v := range c.retentions[1:] {
		queriesMap[fmt.Sprintf("%s%s", ContinuousQueryPrefix, v.name)] = true
	}

	for _, v := range res.Results[0].Series[index].Values {
		if queriesMap[v[columnsMap["name"]].(string)] == false {
			logrus.Warnf("Found unknown Continuous Query in influxdb: '%s'", v[columnsMap["name"]].(string))
		}
		delete(queriesMap, v[columnsMap["name"]].(string))
	}

	if len(queriesMap) == (len(c.retentions) - 1) {
		return false, nil
	}

	if len(queriesMap) > 0 {
		return false, errors.New("some continuous queries are missing from influxdb")
	}
	return true, nil
}

func (c *client) DatabaseExists() (bool, error) {
	query := influx_client.NewQuery("SHOW DATABASES", "", "")
	res, err := c.client.Query(query)
	// path: results[0].series[0].values[0][0]
	if err != nil {
		return false, err
	}
	dbs := res.Results[0].Series[0].Values

	for _, v := range dbs {
		if v[0] == c.db {
			return true, nil
		}
	}
	return false, nil
}

func (c *client) CreateDatabase() error {
	query := influx_client.NewQuery(fmt.Sprintf("CREATE DATABASE \"%s\"", c.db), "", "")
	res, err := c.client.Query(query)
	if err != nil {
		return err
	}
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

// InitDB Checks for database, retention policy and continuous query and creates them all
func InitDB(d databaseManager) error {

	//Database
	exists, err := d.DatabaseExists()
	if err != nil {
		return err
	}
	if !exists {
		logrus.Warn("Creating influxdb database")
		err = d.CreateDatabase()
		if err != nil {
			return err
		}
	}

	// Retention policies
	rps, err := d.RetentionPoliciesExist()
	if err != nil {
		return err
	}
	if !rps {
		logrus.Warn("Creating influxdb retention policies")
		err := d.CreateRetentionPolicies()
		if err != nil {
			return err
		}
	}

	// Continuous queries
	cqs, err := d.ContinuousQueriesExist()
	if err != nil {
		return err
	}
	if !cqs {
		logrus.Warn("Creating influxdb downsampling continuous queries")
		err := d.CreateContinuousQueries()
		if err != nil {
			return err
		}
	}
	return nil
}

// Translate row columns to ints
func columnsAsMap(row models.Row) map[string]int {
	columnsMap := make(map[string]int)

	// Get column names
	for i, v := range row.Columns {
		columnsMap[v] = i
	}
	return columnsMap
}

// getRetentionPolicy returns best suited retention policy for given interval
func (c *client) getRetentionPolicy(begin time.Time, end time.Time) (*retention, error) {
	seconds := (time.Since(begin) - time.Since(end)).Seconds()

	if len(c.retentions) == 0 {
		return nil, errors.New("no retention policies defined")
	}

	for _, v := range c.retentions {
		if v.duration.Seconds() >= seconds {
			return &v, nil
		}
	}
	return &c.retentions[0], errors.New("invalid time range")
}
