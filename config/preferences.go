package config

import (
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

type Preferences struct {
	ServerVersion string
	Host          string
	TokenExpires  bool
	TokenDuration time.Duration
	SecretKey     string
	Logging       LoggingPreferences
	Metrics       bool
}

type LoggingPreferences struct {
	Directory      string
	MainLogFile    string
	LogLevel       logrus.Level
	SqlLogFile     string
	LogSql         bool
	LogRequestFile string
	LogRequests    bool
	LogRotate      bool
	LogMaxFileSize int
	LogFormat      string
}

func LoggingPreferencesFromConfig(l *Logging) (LoggingPreferences, error) {
	lp := LoggingPreferences{
		Directory:   l.Directory,
		LogSql:      l.LogSql,
		LogRequests: l.LogRequests,
	}

	mainlevel, err := logrus.ParseLevel(l.LogLevel)
	if err != nil {
		return lp, err
	}
	logrus.SetLevel(mainlevel)

	lp.MainLogFile = path.Join(l.Directory, l.MainLogFile)
	lp.SqlLogFile = path.Join(l.Directory, l.SqlLogFile)
	lp.LogRequestFile = path.Join(l.Directory, l.LogRequestFile)
	return lp, nil
}
