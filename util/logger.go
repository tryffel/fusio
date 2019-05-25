package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type StdLogger struct {
	Debug bool
}

func (s *StdLogger) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		_, err = os.Stderr.Write([]byte("Unable to read logrus log line"))
	} else {
		_, err = os.Stdout.Write([]byte(line))
	}
	return err
}

// Write these levels to stdout
func (s *StdLogger) Levels() []logrus.Level {
	if s.Debug {
		return []logrus.Level{
			logrus.FatalLevel,
			logrus.PanicLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		}
	} else {

		return []logrus.Level{
			logrus.FatalLevel,
			logrus.PanicLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		}
	}
}

func (s *StdLogger) SetDebug(mode bool) {
	s.Debug = mode
}

type SqlLogger struct {
	Logger *logrus.Logger
}

func (l *SqlLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		query := fmt.Sprintf("'%s'", v[3])
		took := fmt.Sprintf("%0.2f", v[2].(time.Duration).Seconds()*1000)
		returned := v[5]
		l.Logger.WithFields(logrus.Fields{"duration (ms)": took, "type": "sql", "rows": returned}).Print(query)
	}
	if v[0] == "log" {
		l.Logger.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

func (l *SqlLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *SqlLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
