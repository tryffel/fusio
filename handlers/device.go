package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/dtos"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
	"net/http"
	"time"
)

// GetDeviceById Query devices by id
func (h *Handler) GetDeviceById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if !util.IsUuid(id) {
		JsonErrorResponse(w, ResponseInvalidId, http.StatusBadRequest)
		return
	}

	device, err := h.Store.Device.GetById(id)
	if err != nil {
		logrus.Error(err)
	}

	if device.ID == "" {
		JsonErrorResponse(w, ResponseNotFound, http.StatusNotFound)
		return
	}
	err = h.Store.Device.LoadGroups(device)
	if err != nil {
		logrus.Error("", err)
		JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, dtos.FromDevice(device))
}

// GetDevices Get all devices
func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {

	if !h.UserAuthenticated(r) {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}
	userId := r.Context().Value("UserId")
	user, err := h.Store.User.FindByName(userId.(string))
	if err != nil {
		JsonErrorResponse(w, h.Store.Errors.GetUserFriendlyError(err, "user").Error(), http.StatusBadRequest)
		return
	}
	devices, err := h.Store.Device.GetByOwnerId(user.ID)

	if err != nil {
		JsonErrorResponse(w, h.Store.Errors.GetUserFriendlyError(err, "devices").Error(), http.StatusBadRequest)
		return
	}

	dto := make([]dtos.Device, len(*devices))
	for i, v := range *devices {
		err := h.Store.Device.LoadGroups(&v)
		if err != nil {
			logrus.Error(err)
			return
		}
		dto[i] = *dtos.FromDevice(&v)
	}
	JsonResponse(w, dto)
	return

}

// AddDevice Add new device and return api-key
func (h *Handler) AddDevice(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusUnauthorized)
		return
	}

	dto := &dtos.NewDevice{}
	err = dtos.Validate(w, r, dto)
	if err != nil {
		return
	}

	deviceType, err := models.DeviceTypeFromString(dto.Type)
	if err != nil {
		JsonErrorResponse(w, "Invalid device type", http.StatusBadRequest)
		return
	}

	device := &models.Device{
		OwnerId:    user.ID,
		Name:       dto.Name,
		Info:       dto.Info,
		DeviceType: deviceType,
	}

	if device.Name == "" {
		JsonErrorResponse(w, "Device must have a name", http.StatusBadRequest)
		return
	}

	err = h.Store.Device.Create(device)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, device.Name)
		JsonErrorResponse(w, friendly.Error(), http.StatusInternalServerError)
		return
	}
	key, err := h.Store.ApiKey.New("Autocreated", device, h.Preferences.TokenExpires, time.Now().Add(h.Preferences.TokenDuration))

	if err != nil {
		logrus.Error(err)
		JsonErrorResponse(w, "Failed to create device", http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	data["id"] = device.ID
	data["api_key"] = key.Key

	JsonResponse(w, data)
	return

}

// GetMeasurements sends measurements for given device
func (h *Handler) GetDeviceMeasurements(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if user == nil {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "device")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]
	access, err := h.Store.Device.UserHasAccess(user.ID, []string{id})
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "device")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}
	if !access {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}

	measurements, err := h.Store.Measurement.GetDeviceMeasurements(id)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "measurement")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	data := dtos.MeasurementList{
		Measurements: measurements,
	}

	JsonResponse(w, data)
}
