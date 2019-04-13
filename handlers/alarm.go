package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/dtos"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
	"net/http"
	"strings"
	"time"
)

const (
	// How many points to evaluate per interval
	// e.g. for 30 minute interval group by time is 3 min
	AlarmEvaluteOverIntervalNum int = 10
)

type AlarmDTO struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Info        string            `json:"info"`
	Message     string            `json:"message"`
	Fired       bool              `json:"fired"`
	Enabled     bool              `json:"enabled"`
	Group       string            `json:"group"`
	Interval    Interval          `json:"past"`
	Filter      string            `json:"filter"`
	History     []AlarmHistoryDto `json:"history"`
	HistorySize int               `json:"history_size"`
}

type AlarmHistoryDto struct {
	FiredAt   time.Time `json:"fired_at"`
	Value     string    `json:"value"`
	Cleared   bool      `json:"cleared"`
	ClearedAt time.Time `json:"cleared_at"`
}

func AlarmHistoryToDto(a *models.AlarmHistory) *AlarmHistoryDto {
	dto := &AlarmHistoryDto{
		FiredAt: a.FiredAt,
		Value:   a.Value,
		Cleared: a.Cleared,
	}
	if a.Cleared {
		dto.ClearedAt = a.ClearedAt
	}
	return dto
}

func AlarmHistoryArrayToDto(arr *[]models.AlarmHistory) *[]AlarmHistoryDto {

	dto := make([]AlarmHistoryDto, len(*arr))
	for i, v := range *arr {
		dto[i].FiredAt = v.FiredAt
		dto[i].Value = v.Value
		dto[i].Cleared = v.Cleared
		if v.Cleared {
			dto[i].ClearedAt = v.ClearedAt
		}
	}
	return &dto
}

type Inputs map[string]string
type Interval time.Duration

func (i *Interval) UnmarshalJSON(data []byte) error {
	result, err := time.ParseDuration(strings.Trim(string(data), "\""))
	*i = Interval(result)
	return err
}

func (i *Interval) ToDuration() time.Duration {
	return time.Duration(*i)
}

// Divide time interval
func (i *Interval) Divide(num int64) time.Duration {
	total := int64(i.ToDuration()) / num
	return time.Duration(total)
}

func (i *Interval) Seconds() float64 {
	return i.ToDuration().Seconds()
}

// AlarmQueryDto Dto to encapsulate alarm query
type AlarmQueryDto struct {
	Group    string            `json:"group"`
	Inputs   map[string]string `json:"inputs"`
	ForPast  Interval          `json:"past"`
	Filter   string            `json:"filter"`
	Interval Interval          `json:"interval"`
}

func AlarmToDto(a *models.Alarm, history_count int) *AlarmDTO {

	dto := &AlarmDTO{
		ID:          a.ID,
		Name:        a.Name,
		Info:        a.Info,
		Message:     a.Message,
		Fired:       a.Fired,
		Enabled:     a.Enabled,
		Group:       a.Group,
		Filter:      a.Filter.Expression,
		Interval:    Interval(a.RunInterval),
		History:     *AlarmHistoryArrayToDto(&a.History),
		HistorySize: history_count,
	}
	return dto
}

func (h *Handler) GetAlarms(w http.ResponseWriter, r *http.Request) {
	if !h.UserAuthenticated(r) {
		JsonResponse(w, ResponseUnauthorized)
		return
	}

	user, err := h.Store.User.FindByName(r.Context().Value("UserId").(string))
	if err != nil {
		logrus.Error(err, r.Context().Value("UserId"))
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}
	if user.ID < 0 {
		JsonResponse(w, ResponseNotFound)
		logrus.Error("User authenticated but user not found from db: ", r.Context().Value("UserId"))
		return
	}
	alarms, err := h.Store.Alarm.FindByOwner(int(user.ID))
	dto := make([]AlarmDTO, 0)

	for _, v := range *alarms {
		dto = append(dto, *AlarmToDto(&v, -1))
	}

	JsonResponse(w, dto)
}

func (h *Handler) GetAlarmById(w http.ResponseWriter, r *http.Request) {
	if !h.UserAuthenticated(r) {
		JsonResponse(w, ResponseUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	if !util.IsUuid(id) {
		JsonResponse(w, ResponseInvalidId)
		return
	}

	alarm, err := h.Store.Alarm.FindById(id)
	if err != nil {
		logrus.Errorf("Failed to retrieve alarm by id: %s", err)
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}

	err = h.Store.Alarm.LoadHistory(alarm)
	if err != nil {
		logrus.Error(err)
	}

	if alarm.ID == "" {
		JsonResponse(w, ResponseNotFound)
		return
	}

	history, err := h.Store.Alarm.GetHistorySize(alarm)
	if err != nil {
		logrus.Error(err)
	}

	dto := AlarmToDto(alarm, history)
	JsonResponse(w, dto)
	return
}

// CreateAlarm Create new alarm
// Only allow users (not devices) to create alarms
func (h *Handler) CreateAlarm(w http.ResponseWriter, r *http.Request) {
	if !h.UserAuthenticated(r) {
		JsonErrorResponse(w, ResponseAccessDenied, http.StatusForbidden)
		return
	}
	userName := r.Context().Value("UserId")

	user, err := h.Store.User.FindByName(userName.(string))
	if err != nil {
		logrus.Error(err, userName)
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}
	if user.ID < 1 {
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}

	dto := &dtos.NewAlarm{}
	err = dtos.Validate(w, r, dto)
	if err != nil {
		JsonErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	alarm, err := dto.ToAlarm()
	alarm.OwnerId = user.ID
	if err != nil {
		JsonErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Store.Alarm.Create(alarm)
	if err != nil {
		logrus.Error(err)
		JsonErrorResponse(w, ResponseInternalError, http.StatusBadRequest)
		return
	}
	JsonMessage(w, "alarm", alarm.ID)
}
