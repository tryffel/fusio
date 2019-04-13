package handlers

import (
	"github.com/tryffel/fusio/dtos"
	"github.com/tryffel/fusio/notifications"
	"net/http"
)

func (h *Handler) CreateOutputChannel(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}

	// TODO: make validator to validate notification data
	dto := &dtos.NewOutputChannel{}
	err = dtos.Validate(w, r, dto)
	if err != nil {
		return
	}

	channel, err := dto.ToOutputChannel()
	if err != nil {
		JsonErrorResponse(w, ResponseInvalidBody, http.StatusBadRequest)
		return
	}

	err = notifications.ValidateChannel(dto.Type, channel.Data)
	if err != nil {
		JsonErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	channel.OwnerId = user.ID

	err = h.Store.OutputChannel.Create(channel)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "output channel")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	data["id"] = channel.ID
	JsonResponse(w, data)
	return
}

func (h *Handler) CreateOutput(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		JsonErrorResponse(w, ResponseUnauthorized, http.StatusForbidden)
		return
	}

	dto := &dtos.NewOutput{}
	err = dtos.Validate(w, r, dto)

	output := dto.ToOutput()
	output.OwnerId = user.ID

	alarm, err := h.Store.Alarm.FindByOwnerAndId(dto.AlarmId, int(user.ID))
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "alarm")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	if alarm.ID != dto.AlarmId {
		JsonErrorResponse(w, "alarm not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		JsonErrorResponse(w, "invalid template", http.StatusBadRequest)
	}

	err = h.Store.Output.Create(output)
	if err != nil {
		friendly := h.Store.Errors.GetUserFriendlyError(err, "output")
		JsonErrorResponse(w, friendly.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	data["id"] = output.ID
	JsonResponse(w, data)
}
