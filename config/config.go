package config

import (
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	ServiceVersion string = "0.2.0"
	SecretKeySize  int    = 60
)

// Config root part
type Config struct {
	General        General
	Server         Server
	Database       Database
	Influxdb       Influxdb
	Alarms         Alarms
	Metrics        Metrics
	Logging        Logging
	serviceVersion string
	configPath     string
	configFile     string
	preferences    Preferences
}

type Server struct {
	ListenTo string `yaml:"listen_to"`
	Port     int    `yaml:"port"`
}

type General struct {
	Mode                    string `yaml:"mode"`
	RemoveData              bool   `yaml:"remove_old_data"`
	SecretKey               string `yaml:"secret_key"`
	TokenExpires            string `yaml:"token_expire_time"`
	tokenExpiration         time.Duration
	Metrics                 bool `yaml:"metrics"`
	backgroundTasksInterval time.Duration
}

type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	File     string `yaml:"file"`
	Ssl      bool   `yaml:"ssl"`
	loadDemo bool
}

type Influxdb struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type Logging struct {
	Directory      string `yaml:"directory"`
	MainLogFile    string `yaml:"main_file"`
	LogLevel       string `yaml:"level"`
	SqlLogFile     string `yaml:"sql_file"`
	LogSql         bool   `yaml:"log_sql"`
	LogRequestFile string `yaml:"requests_file"`
	LogRequests    bool   `yaml:"log_requests"`
}

type Alarms struct {
	RunBackground bool          `yaml:"enabled"`
	Interval      util.Interval `yaml:"interval"`
}

type Metrics struct {
	RunMetrics bool          `yaml:"enabled"`
	Interval   util.Interval `yaml:"interval"`
}

// Create new configuration
func NewConfig(location string, name string) Config {
	config := Config{
		configPath:     location,
		configFile:     name,
		serviceVersion: ServiceVersion,
	}
	return config
}

// LoadConfig Read configuration file
func (c *Config) LoadConfig() {
	location := filepath.Join(c.configPath, c.configFile)
	data, err := ioutil.ReadFile(location)
	if err != nil {
		logrus.Error("Failed to open configuration file: ", location)
		panic(err)
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		logrus.Fatal("Failed to parse configuration file: ", err)
		panic(err)
	}

	if c.General.SecretKey == "" {
		c.General.SecretKey = util.RandomKey(SecretKeySize)
	}
	c.General.TokenExpires = "14d"
	c.General.tokenExpiration = time.Hour * 24 * 14

	c.preferences.Host = "Fusio"
	c.preferences.ServerVersion = ServiceVersion
	c.preferences.TokenExpires = true
	c.preferences.TokenDuration = c.General.tokenExpiration
	c.preferences.SecretKey = c.General.SecretKey

	c.preferences.Logging, err = LoggingPreferencesFromConfig(&c.Logging)
	if err != nil {
		logrus.Error("Failed to read logging config: ", err)
	}
	c.preferences.Metrics = c.General.Metrics

	c.SaveFile()
}

// CreatFile Create empty configuration file
func (c *Config) CreateFile() {
	c.setDefaultValues()
	c.SaveFile()
}

// SaveFile save file
func (c *Config) SaveFile() {
	location := filepath.Join(c.configPath, c.configFile)
	data, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(location, data, 0600)
	if err != nil {
		logrus.Error("Failed to create empty configuration file: ", location)
		panic(err)
	}
}

// Set default values with sane configuration
func (c *Config) setDefaultValues() {
	c.General.Mode = "prod"

	c.Server.Port = 8080
	c.Server.ListenTo = "0.0.0.0"
}

func (c *Config) GetServerVersion() string {
	return c.serviceVersion
}

func (c *Config) LoadDemoData() {
	c.Database.loadDemo = true
}

func (d *Database) LoadDemoData() bool {
	return d.loadDemo
}

func (c *Config) GetPreferences() *Preferences {
	return &c.preferences
}
