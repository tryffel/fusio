package repository_tests

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"github.com/tryffel/fusio/storage/repository_impl"
	"os"
)

type Database struct {
	db            *gorm.DB
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
}

func (d *Database) LogSql(log bool) {
	d.db.LogMode(log)
}

func getDatabaseFromArgs() *Database {
	db := &Database{}

	dbType := os.Getenv("fusio_test_db_type")
	host := os.Getenv("fusio_test_db_host")
	port := os.Getenv("fusio_test_db_port")
	user := os.Getenv("fusio_test_db_user")
	password := os.Getenv("fusio_test_db_password")
	database := os.Getenv("fusio_test_db_database")
	file := os.Getenv("fusio_test_db_file")

	if dbType == "" {
		logrus.Error("Failed to read environment variables")
		return nil
	}

	var err error
	var url string

	switch dbType {
	case "postgres":
		if password == "" {
			url = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
				host, port, user, database)
		} else {
			url = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s, sslmode=disable",
				host, port, user, database, password)
		}

		db.db, err = gorm.Open("postgres", url)
		if err != nil {
			logrus.Error(err)
		}
	case "sqlite":
		url := fmt.Sprintf("file=%s", file)
		db.db, err = gorm.Open("sqlite", url)
		if err != nil {
			logrus.Error(err)
		}
	default:
		logrus.Error("Unknown database type")
		return nil
	}

	db.Alarm = repository_impl.NewAlarmRepository(db.db)
	db.Group = repository_impl.NewGroupRepository(db.db)
	//db.Measurement = repository_impl.meas
	db.User = repository_impl.NewUserRepository(db.db)
	db.Device = repository_impl.NewDeviceRepository(db.db)
	db.ApiKey = repository_impl.NewApiKeyRepository(db.db)
	db.Output = repository_impl.NewOutputRepository(db.db)
	db.OutputChannel = repository_impl.NewOutputChannelRepository(db.db)
	db.Errors = repository_impl.NewErrors(dbType)
	db.Pipeline = repository_impl.NewPipelineRepository(db.db)

	db.db.AutoMigrate(&models.User{})
	db.db.AutoMigrate(&models.Device{})
	db.db.AutoMigrate(&models.Group{})
	db.db.AutoMigrate(&models.Alarm{})
	db.db.AutoMigrate(&models.AlarmHistory{})
	db.db.AutoMigrate(&models.ApiKey{})
	db.db.AutoMigrate(&models.OutputChannel{})
	db.db.AutoMigrate(&models.Output{})
	db.db.AutoMigrate(&models.OutputHistory{})

	return db
}

//createUser creates user to use on tests
func (db *Database) createUser() *models.User {
	user := &models.User{
		Name:     "test_user",
		Email:    "test@test.test",
		IsActive: true,
		IsAdmin:  false,
	}

	// Create user for first time
	// Should succeed
	err := db.User.Create(user, "password")
	if err != nil {
		logrus.Errorf("cannot create user: %s", err.Error())
	}
	return user
}

func (db *Database) createGroup(ownerId uint, name string) *models.Group {
	group := &models.Group{
		Name:    name,
		OwnerId: ownerId,
	}

	err := db.Group.Create(group)
	if err != nil {
		logrus.Error("Failed to create group: ", err.Error())
	}
	return group
}

func (db *Database) RemoveAllRecords() {

	db.db.Unscoped().Delete(&models.User{})
	db.db.Unscoped().Delete(&models.Group{})
	db.db.Unscoped().Delete(&models.Group{})
	db.db.Unscoped().Delete(&models.Alarm{})
	db.db.Unscoped().Delete(&models.AlarmHistory{})
	db.db.Unscoped().Delete(&models.ApiKey{})
	db.db.Unscoped().Delete(&models.OutputChannel{})
	db.db.Unscoped().Delete(&models.Output{})
	db.db.Unscoped().Delete(&models.OutputHistory{})
	db.db.Unscoped().Delete(&models.Pipeline{})
}

func getRedisFromArgs() (repository.Cache, error) {
	url := os.Getenv("fusio_test_redis_url")

	if url == "" {
		return nil, errors.New("no redis url defined {fusio_test_redis_url}")
	}

	c := config.Redis{
		Type: "tcp",
		Url:  url,
	}

	return repository_impl.NewRedis(&c)
}
