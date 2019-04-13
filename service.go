package fusio

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/alarm"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/dtos"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/handlers"
	"github.com/tryffel/fusio/metrics"
	"github.com/tryffel/fusio/storage"
	"github.com/tryffel/fusio/util"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type Service struct {
	Config       *config.Config
	Store        *storage.Store
	BaseRouter   *mux.Router
	ApiRouter    *mux.Router
	PublicRouter *mux.Router
	Server       *http.Server
	Handler      handlers.Handler
	AlarmTask    *alarm.BackgroundTask
	MetricsTask  *metrics.BackgroundTask
	lock         sync.RWMutex
	logRequest   *os.File
	logSql       *os.File
	logFile      *os.File
	running      bool
}

func NewService(config *config.Config) (*Service, error) {
	service := &Service{
		Config: config,
	}

	logFormat := &prefixed.TextFormatter{
		FullTimestamp:  true,
		QuoteCharacter: "'",
	}
	logFormat.ForceFormatting = true
	logrus.SetFormatter(logFormat)

	var requestLogger *logrus.Logger
	var sqlLogger *logrus.Logger

	err := service.openLogFiles()
	if err != nil {
		err = Err.Wrap(&err, "Could not open all log files")
		e := service.closeLogFiles()
		if e != nil {
			err = Err.Wrap(&err, e.Error())
		}
		return service, err
	}

	if service.Config.Logging.LogRequests {
		requestLogger = logrus.New()
		requestLogger.SetOutput(service.logRequest)
		requestLogger.SetFormatter(logFormat)
		requestLogger.SetLevel(logrus.InfoLevel)
	}
	if service.Config.Logging.LogSql {
		sqlLogger = logrus.New()
		sqlLogger.SetOutput(service.logSql)
		sqlLogger.SetFormatter(logFormat)
		sqlLogger.SetLevel(logrus.InfoLevel)
	}

	sqlLog := &util.SqlLogger{
		Logger: sqlLogger,
	}

	logrus.SetOutput(service.logFile)
	logrus.AddHook(&util.StdLogger{})

	store, err := storage.NewStore(&config.Database, &config.Influxdb, &config.GetPreferences().Logging, sqlLog)
	if err != nil {
		logrus.Fatal("Failed initializing database connection: ", err)
		panic("Failed to initialize database")
		return service, err
	}
	logrus.Debug("Database connection ok")

	service.Store = store

	service.BaseRouter = mux.NewRouter()
	service.ApiRouter = service.NewApiRouter()
	service.PublicRouter = service.NewApiRouter()
	service.CreateRoutes()

	addr := fmt.Sprintf("%s:%d", service.Config.Server.ListenTo, service.Config.Server.Port)
	service.Server = &http.Server{
		Handler:      service.BaseRouter,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Enable custom dto validators
	dtos.AddUuidArrayValidator()
	dtos.AddIntervalValidator()
	dtos.AddDurationValidator()

	service.MetricsTask, err = metrics.NewBackgroundTask(config, service.Store)
	if err != nil {
		return service, err
	}
	service.AlarmTask, err = alarm.NewBackgroundTask(*config, service.Store, service.MetricsTask)
	if err != nil {
		logrus.Fatal("", err)
		return service, err
	}

	service.Handler = handlers.NewHandler(service.Store, service.MetricsTask, *service.Config.GetPreferences(), requestLogger)
	return service, nil
}

// Start Start service
func (s *Service) Start() {
	s.lock.Lock()
	if !s.running {
		s.running = true
		s.lock.Unlock()

		logrus.Info(fmt.Sprintf("----- Fusio Server %s -----", s.Config.GetServerVersion()))
		logrus.Info("")
		logrus.Warn("Starting service.")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		signal.Notify(c, os.Interrupt, syscall.SIGINT)
		go func() {
			<-c
			s.Stop()
		}()

		if s.Config.Alarms.RunBackground {
			err := s.AlarmTask.Start()
			if err != nil {
				logrus.Errorf("Error starting background alarms: %s", err)
			}
		} else {
			logrus.Info("Alarm task disabled")
		}

		logrus.Info("Listening on ", s.Server.Addr)
		if s.Config.Metrics.RunMetrics {
			err := s.MetricsTask.Start()
			if err != nil {
				logrus.Error("Error starting background alarms: ", err)
			}
		} else {
			logrus.Info("Metrics disabled")
		}

		err := s.Server.ListenAndServe()
		if err != nil {
			if err.Error() != "http: Server closed" {
				logrus.Error(err.Error())
			}
		}
	} else {
		s.lock.Unlock()
	}

}

// Stop Stop service
func (s *Service) Stop() {
	s.lock.Lock()
	if s.running {
		s.running = false
		s.lock.Unlock()
		logrus.Warn("Stopping service")

		s.AlarmTask.Stop()
		s.MetricsTask.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := s.Server.Shutdown(ctx)
		if err != nil {
			if err.Error() != "http: Server closed" {
				logrus.Error(err)
			}
			err := s.Store.Close()
			if err != nil {
				logrus.Error(err)
			}

			err = s.Server.Close()
			if err != nil {
				logrus.Error(err.Error())
			}
		}
		err = s.closeLogFiles()
		if err != nil {
			Err.Log(err)
		}
	}
}

func (s *Service) openLogFiles() error {

	mode := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	perm := os.FileMode(0760)

	conf := &s.Config.Logging
	file, err := os.OpenFile(filepath.Join(conf.Directory, conf.MainLogFile), mode, perm)
	if err != nil {
		return err
	}
	s.logFile = file

	if conf.LogSql {
		file, err = os.OpenFile(filepath.Join(conf.Directory, conf.SqlLogFile), mode, perm)
		if err != nil {
			return err
		}
		s.logSql = file
	}
	if conf.LogRequests {
		file, err = os.OpenFile(filepath.Join(conf.Directory, conf.LogRequestFile), mode, perm)
		if err != nil {
			return err
		}
		s.logRequest = file
	}
	return nil
}

func (s *Service) closeLogFiles() error {
	var e error
	err := s.logFile.Close()
	if err != nil {
		e = Err.Wrap(&err, "failed to close main log file")
	}
	if s.logSql != nil {
		err = s.logSql.Close()
		if err != nil {
			e = Err.Wrap(&e, "failed to close sql log file")
		}
	}
	if s.logRequest != nil {
		err = s.logRequest.Close()
		if err != nil {
			e = Err.Wrap(&e, "failed to close requests log file")
		}
	}
	return e
}
