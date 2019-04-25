package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
)

type Database struct {
	engine      *gorm.DB
	initialized bool
	db          string
	log         *util.SqlLogger
}

// NewDatabase initialize new database
func NewDatabase(c *config.Database, p *config.LoggingPreferences, logger *util.SqlLogger) (*Database, error) {

	db := &Database{}

	var engine *gorm.DB
	var err error

	db.log = logger

	switch c.Type {
	case "sqlite":
		engine, err = gorm.Open("sqlite3", c.File)
		db.db = c.File
	case "postgres":
		url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			c.Host, c.Port, c.Username, c.Database, c.Password)
		engine, err = gorm.Open("postgres", url)
		db.db = fmt.Sprintf("postgres://%s:%s/%s", c.Host, c.Port, c.Database)
	case "mysql":
		engine, err = gorm.Open("mysql",
			fmt.Sprintf("%s:%s@%s:%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
				c.Username, c.Password, c.Host, c.Port, c.Database))
		db.db = fmt.Sprintf("mysql://%s:%s/%s", c.Host, c.Port, c.Database)

	default:
		logrus.Fatal("Invalid database type configured!")
		panic("Invalid database configured")
		return nil, err
	}

	if err != nil {
		return db, err
	}

	db.engine = engine

	if err != nil {
		return db, err
	}
	db.engine.LogMode(p.LogSql)
	if p.LogSql {
		db.engine.SetLogger(db.log)
	}

	db.runMigrations()

	if c.LoadDemoData() {
		logrus.Info("Loading demo data")
		db.DemoData()
	}
	db.initialized = true
	return db, nil
}

// Close closes database connection
func (db *Database) Close() error {
	//db.sqlLogger.Close()
	return db.engine.Close()
}

// GetEngine Get initialized database engine
// TODO: remove method
func (db *Database) GetEngine() *gorm.DB {
	if db.initialized {
		return db.engine
	}
	return nil
}

func (db *Database) runMigrations() {
	logrus.Debug("Running db migrations")

}

func (db *Database) DemoData() {
	logrus.Warn("Filling database with demo data")

	user1 := models.User{
		Name:     "TestUser",
		Email:    "TestEmail@test.com",
		IsActive: true,
		IsAdmin:  true,
	}
	err := user1.SetPassword("12345")
	if err != nil {
		logrus.Error(err.Error())
	}
	db.engine.Create(&user1)

	dev1 := models.Device{
		Name:       "Sensor1",
		OwnerId:    user1.ID,
		Info:       "Demo device",
		DeviceType: models.DeviceSensor,
	}
	db.engine.Create(&dev1)

	dev2 := models.Device{
		Name:       "Sensor2",
		OwnerId:    user1.ID,
		Info:       "Another demo device",
		DeviceType: models.DeviceSensor,
	}
	db.engine.Create(&dev2)

	dev3 := models.Device{
		Name:       "Controller1",
		OwnerId:    user1.ID,
		Info:       "Controller demo device",
		DeviceType: models.DeviceController,
	}
	db.engine.Create(&dev3)

	group1 := models.Group{
		Name:    "TestGroup1",
		Info:    "Group for testing purposes",
		OwnerId: user1.ID,
	}

	db.engine.Create(&group1)
	err = group1.AddDevice(db.engine, &dev1)
	if err != nil {
		logrus.Error(err.Error())
	}
	err = group1.AddDevice(db.engine, &dev2)
	if err != nil {
		logrus.Error(err.Error())
	}
	err = group1.AddDevice(db.engine, &dev3)
	if err != nil {
		logrus.Error(err.Error())
	}

	key1 := models.ApiKey{
		Name:     "Test key 1",
		DeviceId: dev1.ID,
		Key:      "igqPGmSzQz38tozOOfuK",
	}
	db.engine.Create(&key1)

	key2 := models.ApiKey{
		Name:     "Test key 2",
		DeviceId: dev2.ID,
		Key:      "d2b4MtsynTQZnXgYX5N0",
	}
	db.engine.Create(&key2)

	key3 := models.ApiKey{
		Name:     "Controller key",
		DeviceId: dev3.ID,
		Key:      "Kb4gmkuXfxYAZSXJ8dpG",
	}
	db.engine.Create(&key3)
}
