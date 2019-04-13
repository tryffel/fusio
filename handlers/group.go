package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/dtos"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"github.com/tryffel/fusio/util"
	"net/http"
)

// GroupDTO Group to expose on api
type GroupDTO struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Info    string          `json:"info"`
	Devices []DeviceInfoDTO `json:"devices"`
}

// GroupDTOFromGroup Create groupDto from models.Group
func GroupDTOFromGroup(g *models.Group) *GroupDTO {
	group := &GroupDTO{
		Id:   g.ID,
		Name: g.Name,
		Info: g.Info,
	}

	for i := range g.Devices {
		group.Devices = append(group.Devices, *DeviceInfoDTOFromDevice(&g.Devices[i]))
	}

	return group
}

// DeviceInfoDTO Device to expose on api
type DeviceInfoDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// DeviceInfoDTOFromDevice Create deviceDTO from models.Device
func DeviceInfoDTOFromDevice(d *models.Device) *DeviceInfoDTO {
	device := &DeviceInfoDTO{
		Id:   d.ID,
		Name: d.Name,
		Type: d.DeviceType.ToString(),
	}
	return device
}

// GetGroupById Get groupDto by id
func (h *Handler) GetGroupById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if util.IsUuid(id) == false {
		JsonErrorResponse(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	group := &models.Group{}
	h.Store.GetDb().Where("id = ?", id).First(&group)
	if group.ID == "" {
		JsonErrorResponse(w, "Group not found", http.StatusNotFound)
		return
	}
	err := group.LoadDevices(h.Store.GetDb())
	if err != nil {
		logrus.Error("Failed loading group devices: ", err)
		JsonErrorResponse(w, "Internal error", http.StatusInternalServerError)
		return
	}

	JsonResponse(w, dtos.GroupToDto(group))
}

// GetGroups Get all groups user has access to
func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	if !h.UserAuthenticated(r) {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}

	user, err := h.Store.User.FindByName(r.Context().Value("UserId").(string))
	if err != nil {
		logrus.Error(err)
		return
	}
	groups, err := h.Store.Group.FindByOwner(int(user.ID))
	if err != nil {
		logrus.Error(err)
		return
	}

	dto := make([]dtos.Group, len(*groups))
	for i, v := range *groups {
		err := v.LoadDevices(h.Store.GetDb())
		if err != nil {
			logrus.Error("", err)
		}
		dto[i] = *dtos.GroupToDto(&v)
	}
	JsonResponse(w, dto)
	return
}

// CreateGroup Create new group
// Post body is groupdto with optional id, name and info. Also can contain devices with ids
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if !h.UserAuthenticated(r) {
		JsonResponse(w, ResponseUnauthorized)
		return
	}

	dto := dtos.NewGroup{}
	dto.Devices = []string{}
	err := dtos.Validate(w, r, &dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	group := dto.ToGroup()

	// Check that devices exist
	if len(group.Devices) > 0 {
		devices := &[]models.Device{}
		existing := 0
		res := h.Store.GetDb().Where("id IN (?)", *group.GetDeviceIds()).Find(&devices).Count(&existing)
		if res.Error != nil {
			logrus.Error(err)
		}
		if len(group.Devices) != existing {
			JsonErrorResponse(w, "Some devices don't exist", http.StatusNotFound)
			return
		}
	}

	h.Store.GetDb().Create(&group)

	JsonMessage(w, "group", group.ID)
}

func (h *Handler) AddDeviceToGroup(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "user")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}
	if user == nil {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}

	groupId := mux.Vars(r)["id"]
	if !util.IsUuid(groupId) {
		JsonErrorResponse(w, ResponseInvalidId, http.StatusBadRequest)
		return
	}

	dto := dtos.IdList{}
	err = dtos.Validate(w, r, &dto)
	if err != nil {
		return
	}

	group, err := h.Store.Group.FindByOwnerAndId(user.ID, groupId)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "group")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}
	if group.ID != groupId {
		JsonErrorResponse(w, "Group not found", http.StatusBadRequest)
		return
	}

	err = h.Store.Group.AddDevices(group, dto.Ids)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "device")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	JsonMessage(w, "Status", "Ok")
}

func (h *Handler) GetGroupMeasurements(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if user == nil {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "group")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]
	access, err := h.Store.Group.UserHasAccess(user.ID, []string{id})
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "group")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}
	if !access {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}

	measurements, err := h.Store.Measurement.GetGroupMeasurements(id)
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

func (h *Handler) SearchGroupByName(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		code := Err.GetErrCode(err)
		if code == Err.Econflict {
			JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
			return
		}
	}

	name := r.URL.Query().Get("name")

	groups, err := h.Store.Group.SearchGroups(&repository.SearchGroupsOpts{OwnerId: user.ID, Name: name})
	if len(*groups) > 0 {
		logrus.Info("Found groups yay")
	}
	dto := make([]dtos.Group, len(*groups))
	for i, v := range *groups {
		err := v.LoadDevices(h.Store.GetDb())
		if err != nil {
			logrus.Error("", err)
		}
		dto[i] = *dtos.GroupToDto(&v)
	}
	JsonResponse(w, dto)
	return
}

func (h *Handler) GetGroupDevices(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		code := Err.GetErrCode(err)
		if code == Err.Econflict {
			JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
			return
		}
	}

	group := mux.Vars(r)["id"]
	devices, err := h.Store.Group.GetDevices(user.ID, group)
	if err != nil {
		e := Err.Wrap(&err, "Failed to get group devices")
		Err.Log(e)
		JsonErrorResponse(w, e.Message, http.StatusBadRequest)
		return
	}
	logrus.Infof("Found %d devices", len(*devices))
	JsonResponse(w, dtos.IdList{Ids: *devices})
}
