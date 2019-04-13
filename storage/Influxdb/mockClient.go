package Influxdb

import (
	"errors"
	influx_client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"time"
)

// TODO: implement influxdb batchpoints to store measurements

// Influxdb mock client
type mockClient struct {
	LastCommand   string
	LastRPCommand string
	LastCQCommand string
	DbExists      bool
	RPsExists     bool
	CQsExists     bool
	StoreData     bool
	Metrics       map[string][]Point
}

func (m *mockClient) CreateRetentionPolicies() error {
	// TODO validate rp query
	return nil
}

func (m *mockClient) CreateContinuousQueries() error {
	// TODO: validate cq query
	return nil
}

func (m *mockClient) RetentionPoliciesExist() (bool, error) {
	return m.RPsExists, nil
}

func (m *mockClient) ContinuousQueriesExist() (bool, error) {
	return m.CQsExists, nil
}

func (m *mockClient) DatabaseExists() (bool, error) {
	return m.DbExists, nil
}

func (m *mockClient) CreateDatabase() error {
	panic("implement me")
}

func (m *mockClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Second, "1.5", nil
}

func (m *mockClient) Write(bp influx_client.BatchPoints) error {
	if !m.StoreData {
		return nil
	}
	return nil
}

func (m *mockClient) Query(q influx_client.Query) (*influx_client.Response, error) {
	logrus.Info(q.Command)
	return &influx_client.Response{}, nil
}

func (m *mockClient) QueryAsChunk(q influx_client.Query) (*influx_client.ChunkedResponse, error) {
	panic("implement me")
}

func (m *mockClient) Close() error {
	return nil
}

func newMockClient(config *config.Influxdb, metrics bool) (Client, error) {
	c := client{
		client:     &mockClient{StoreData: true},
		db:         config.Database,
		retentions: retentions,
	}

	if len(retentions) == 0 {
		return &c, errors.New("no retention policies defined")
	}

	return &c, nil
}
