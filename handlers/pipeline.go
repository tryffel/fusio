package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/dtos"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/util"
	"net/http"
)

func (h *Handler) GetPipelineById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if util.IsUuid(id) {
		JsonErrorResponse(w, "invalid pipeline id", http.StatusBadRequest)
		return
	}

	user, err := h.getUser(r)
	if err != nil {
		if e, ok := err.(*Err.Error); ok {
			JsonErrorResponse(w, e.EndUserMessage(), http.StatusBadRequest)
		}
		return
	}

	pipeline, err := h.Store.Pipeline.FindbyOwnerAndId(user.ID, id)
	if err != nil {
		str, _ := Err.Wrap(&err, "could not get pipeline").EndUserError()
		JsonErrorResponse(w, str, http.StatusBadRequest)
		return
	}

	JsonResponse(w, pipeline)
}

func (h *Handler) CreatePipeline(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		str, _ := Err.Wrap(&err, "no authentication").EndUserError()
		JsonErrorResponse(w, str, http.StatusForbidden)
		return
	}

	panic("Not implemented")
	dto := &dtos.Pipeline{}

	err = dtos.Validate(w, r, dto)
	if err != nil {
		return
	}

	logrus.Info(dto)

	// TODO: complete
	if user.ID == 0 {
		logrus.Info("No user")
	}

}
