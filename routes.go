package fusio

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	API_V1_ROUTE string = "/api/v1"
)

// CreateRoutes Initialize routes with middlewares
// BaseRouter and ApiRouter has to be created
func (s *Service) CreateRoutes() {

	logrus.Debug("Using recoverMiddlerware")
	s.BaseRouter.Use(s.Handler.RecoverMiddleware)
	s.PublicRouter.Use(s.Handler.RecoverMiddleware)

	if s.Config.Logging.LogRequests {
		logrus.Debug("Using loggingmiddleware")
		s.BaseRouter.Use(s.Handler.LoggingMiddleware)
	}

	s.PublicRouter.HandleFunc("/auth/login", s.Handler.Login).Methods("POST")
	s.PublicRouter.HandleFunc("/server", s.Handler.ServerInfo).Methods("GET")

	logrus.Debug("Using Authenticationmiddleware")
	s.ApiRouter.Use(s.Handler.AuthenticationMiddleware)

	/* MEASUREMENTS */
	s.ApiRouter.HandleFunc("/measurements/{measurement}/{aggregation:last|mean|max|min}", s.Handler.GetMeasurement).Methods("GET")
	s.ApiRouter.HandleFunc("/measurements", s.Handler.PutMeasurement).Methods("POST")

	/* GROUPS */
	s.ApiRouter.HandleFunc("/groups", s.Handler.GetGroups).Methods("GET")
	s.ApiRouter.HandleFunc("/groups/{id:[0-9a-zA-Z-]{36}}", s.Handler.GetGroupById).Methods("GET")
	s.ApiRouter.HandleFunc("/groups", s.Handler.CreateGroup).Methods("POST")
	s.ApiRouter.HandleFunc("/groups/{id}/devices", s.Handler.AddDeviceToGroup).Methods("POST")
	s.ApiRouter.HandleFunc("/groups/{id:[0-9a-zA-Z-]{36}}/devices", s.Handler.GetGroupDevices).Methods("GET")
	s.ApiRouter.HandleFunc("/groups/{id}/measurements", s.Handler.GetGroupMeasurements).Methods("GET")

	/* DEVICES */
	s.ApiRouter.HandleFunc("/devices", s.Handler.GetDevices).Methods("GET")
	s.ApiRouter.HandleFunc("/devices/{id}", s.Handler.GetDeviceById).Methods("GET")
	s.ApiRouter.HandleFunc("/devices", s.Handler.AddDevice).Methods("POST")
	s.ApiRouter.HandleFunc("/devices/{id}/measurements", s.Handler.GetDeviceMeasurements).Methods("GET")

	/* ALARMS */
	s.ApiRouter.HandleFunc("/alarms", s.Handler.CreateAlarm).Methods("POST")
	s.ApiRouter.HandleFunc("/alarms", s.Handler.GetAlarms).Methods("GET")
	s.ApiRouter.HandleFunc("/alarms/alarm/{id}", s.Handler.GetAlarmById).Methods("GET")

	/* OUTPUTS */
	s.ApiRouter.HandleFunc("/alarms/outputs", s.Handler.CreateOutput).Methods("POST")
	s.ApiRouter.HandleFunc("/alarms/outputchannels", s.Handler.CreateOutputChannel).Methods("POST")

	/* PIPELINES */
	s.ApiRouter.HandleFunc("/pipelines/{id}", s.Handler.GetPipelineById).Methods("GET")
	s.ApiRouter.HandleFunc("/pipelines", s.Handler.CreatePipeline).Methods("POST")

	/* SEARCH */
	s.ApiRouter.HandleFunc("/groups/search", s.Handler.SearchGroupByName).Methods("GET")

}

func (s *Service) NewApiRouter() *mux.Router {
	return s.BaseRouter.PathPrefix(API_V1_ROUTE).Subrouter()
}
