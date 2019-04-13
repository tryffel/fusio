package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/util"
	"net/http"
)

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"passowrd"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	dto := &LoginDto{}
	err := DtoFromRequest(r, &dto)
	if err != nil {
		JsonErrorResponse(w, ResponseInvalidBody, http.StatusBadRequest)
		return
	}

	user, err := h.Store.User.FindByName(dto.Username)
	if err != nil {
		JsonErrorResponse(w, ResponsePasswordFailed, http.StatusForbidden)
		return
	}

	if user.Name == "" {
		JsonErrorResponse(w, ResponsePasswordFailed, http.StatusForbidden)
		return
	}

	if util.PasswordMatches(dto.Password, user.Password) {
		token, err := util.NewToken(user.LowerName, h.Preferences.SecretKey, h.Preferences.TokenDuration)
		if err != nil {
			logrus.Error("Failed to issue token: ", err)
			JsonErrorResponse(w, ResponseInternalError, http.StatusInternalServerError)
			return
		}
		JsonMessage(w, "token", token)
		return
	}
	JsonErrorResponse(w, ResponsePasswordFailed, http.StatusBadRequest)
	return
}
