package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/storage/Influxdb"
	"github.com/tryffel/fusio/storage/repository"
	"github.com/tryffel/fusio/storage/repository_impl"
	"github.com/tryffel/fusio/storage/repository_mock"
	"github.com/tryffel/fusio/util"
)

type Store struct {
	database *Database
	engine   string
	influxdb Influxdb.Client

	Alarm         repository.Alarm
	Device        repository.Device
	Group         repository.Group
	Measurement   repository.Measurement
	Setting       repository.Setting
	User          repository.User
	ApiKey        repository.ApiKey
	Output        repository.Output
	OutputChannel repository.OutputChannel
	Errors        repository.Errors
	Pipeline      repository.Pipeline
	PipelineBlock repository.PipelineBlock
	Cache         repository.Cache
}

func NewStore(conf *config.Config, confInflux *config.Influxdb, logging *config.LoggingPreferences, sqlLogger *util.SqlLogger) (*Store, error) {
	store := &Store{}
	confDb := &conf.Database

	database, err := NewDatabase(confDb, logging, sqlLogger)
	if err != nil {
		return store, err
	}
	store.database = database

	store.engine = confDb.Type

	influx, err := Influxdb.NewClient(confInflux, true, logging.LogSql, sqlLogger)
	if err != nil {
		return store, err
	}
	store.influxdb = influx

	store.Alarm = repository_impl.NewAlarmRepository(store.database.GetEngine())
	store.Group = repository_impl.NewGroupRepository(store.database.GetEngine())
	store.Measurement = repository_impl.NewMeasurementRepository(store.database.GetEngine(), store.influxdb)
	store.User = repository_impl.NewUserRepository(store.database.GetEngine())
	store.Device = repository_impl.NewDeviceRepository(store.database.GetEngine())
	store.ApiKey = repository_impl.NewApiKeyRepository(store.database.GetEngine())
	store.Output = repository_impl.NewOutputRepository(store.database.GetEngine())
	store.OutputChannel = repository_impl.NewOutputChannelRepository(store.database.GetEngine())
	store.Errors = repository_impl.NewErrors(store.engine)
	store.Cache, err = repository_impl.NewRedis(&conf.Redis)
	return store, err
}

func (s *Store) Close() error {
	return s.database.Close()
}

func (s *Store) GetDb() *gorm.DB {
	return s.database.GetEngine()
}

func NewMockStore() (*Store, error) {
	store := &Store{}
	store.Alarm = &repository_mock.MockAlarmRepository{}
	store.Measurement = repository_mock.NewMockMeasurementRepository()
	return store, nil
}
