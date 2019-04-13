package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/metrics"
	"github.com/tryffel/fusio/storage"
	"net/http"
)

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (h *Handler) ServerInfo(w http.ResponseWriter, r *http.Request) {
	info := &ServerInfo{
		Name:    h.Preferences.Host,
		Version: h.Preferences.ServerVersion,
	}
	JsonResponse(w, info)
	return
}

// Handler handler struct that has db connections and handler methods
type Handler struct {
	Store       *storage.Store
	Preferences *config.Preferences
	Metrics     metrics.Metrics
	RequestsLog *logrus.Logger
}

// NewHandler Create new handler
func NewHandler(store *storage.Store, metrics metrics.Metrics, pref config.Preferences, logger *logrus.Logger) Handler {
	h := Handler{
		Store:       store,
		Preferences: &pref,
		Metrics:     metrics,
		RequestsLog: logger,
	}
	return h
}

// Check if user has authenticated
func (h *Handler) UserAuthenticated(r *http.Request) bool {
	user := r.Context().Value("UserId")
	if user == "" || user == nil {
		return false
	}
	return true
}

// Check if device is authenticated
func (h *Handler) DeviceAuthenticated(r *http.Request) bool {
	device := r.Context().Value("DeviceId")
	if device == "" || device == nil {
		return false
	}
	return true
}
