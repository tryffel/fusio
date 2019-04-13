package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/storage/Influxdb"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetMeasurement(w http.ResponseWriter, r *http.Request) {

	deviceId := r.Context().Value("DeviceId")
	measurementName := mux.Vars(r)["measurement"]
	device, err := h.Store.Device.GetById(deviceId.(string))
	if err != nil {
		logrus.Error(err)
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}
	err = h.Store.Device.LoadGroups(device)
	if err != nil {
		logrus.Error(err)
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}

	aggregation := mux.Vars(r)["aggregation"]
	timeRange := r.URL.Query().Get("range")
	n := r.URL.Query().Get("n")
	if timeRange == "" {
		timeRange = "1h"
	}
	if n == "" {
		n = "20"
	}

	duration, err := time.ParseDuration(timeRange)
	if err != nil {
		JsonErrorResponse(w, fmt.Sprintf("Invalid range: %s", timeRange), http.StatusBadRequest)
		return
	}

	number, err := strconv.ParseInt(n, 10, 32)
	if err != nil {
		JsonErrorResponse(w, fmt.Sprintf("Invalid number: %s", n), http.StatusBadRequest)
		return
	}

	if aggregation == "last" {
		number = 1
	}

	filter, err := Influxdb.FilterFromString(fmt.Sprintf("%s(%s)", aggregation, measurementName))
	if err != nil {
		JsonErrorResponse(w, ResponseInvalidBody, http.StatusBadRequest)
		return
	}

	err = h.Store.Device.LoadGroups(device)
	if err != nil {
		logrus.Error(err)
		return
	}

	m, err := h.Store.Measurement.Read(device.ID, "", *filter, time.Now().Add(-duration), time.Now(), number)
	if err != nil {
		JsonErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JsonResponse(w, m)
}

func (h *Handler) PutMeasurement(w http.ResponseWriter, r *http.Request) {

	deviceId := r.Context().Value("DeviceId")

	device, err := h.Store.Device.GetById(deviceId.(string))
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "device")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		h.Metrics.CounterIncrease("http_measurement_insert_fail", 1)
		return
	}
	err = h.Store.Device.LoadGroups(device)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "device")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		h.Metrics.CounterIncrease("http_measurement_insert_fail", 1)
		return
	}

	if err != nil {
		friendlyErr := h.Store.Errors.GetUserFriendlyError(err, "device")
		JsonErrorResponse(w, friendlyErr.Error(), http.StatusBadRequest)
		h.Metrics.CounterIncrease("http_measurement_insert_fail", 1)
		return
	}

	data := make(map[string]float32)
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		JsonErrorResponse(w, ResponseInvalidJson, http.StatusBadRequest)
		h.Metrics.CounterIncrease("http_measurement_insert_fail", 1)
		return
	}
	h.Metrics.CounterIncrease("http_measurement_insert_success", 1)
	timestamp := time.Now()

	measurements := map[string]Influxdb.Point{}

	h.Metrics.CounterIncrease("http_measurement_insert", float64(len(data)))

	for k := range data {
		measurements[k] = Influxdb.Point{
			Timestamp: timestamp,
			Value:     data[k],
		}

	}
	err = h.Store.Measurement.Write(device, measurements)
	if err != nil {
		logrus.Error("", err)
		JsonErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JsonMessage(w, ResponseStatus, ResponseOk)
}
