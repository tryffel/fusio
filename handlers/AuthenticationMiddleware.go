package handlers

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/util"
	"net/http"
	"strings"
)

const (
	HEADER_API_KEY string = "API-KEY"
)

// AuthenticationMiddlware Authenticate users and devices and deny access if unauthorized
// Adds request context with either device or user
// If user, add key "UserId"
// If device, add key "DeviceId"
func (h *Handler) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Device
		key := r.Header.Get(HEADER_API_KEY)
		if key != "" {

			device, err := h.getDevice(key)
			if err == nil {
				ctx := context.WithValue(r.Context(), "DeviceId", device)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// User
		text := r.Header.Get("Authorization")
		parts := strings.Split(text, " ")
		if len(parts) != 2 {
			JsonErrorResponse(w, ResponseInvalidToken, http.StatusForbidden)
			return
		}
		token := parts[1]

		user, err := util.UserFromToken(token, h.Preferences.SecretKey, h.Preferences.TokenExpires,
			h.Preferences.TokenDuration)
		if err != nil {
			JsonErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "UserId", user)
		logrus.Debug("User: ", user)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

// getDevice Get device model from api key
// In future, this could utilize cache, e.g. redis-server to speed up requests
func (h *Handler) getDevice(key string) (string, error) {
	device, err := h.Store.ApiKey.GetDeviceId(key)
	if err != nil {
		return device, errors.New("forbidden")
	}
	return device, nil
}

// getUser Authenticate user against token
func (h *Handler) getUser(r *http.Request) (*models.User, error) {
	userName := r.Context().Value("UserId")

	if userName == nil {
		return nil, &Err.Error{Code: Err.Econflict, Message: "Authentication required",
			Err: errors.New("user not found")}
	}

	if userName.(string) == "" {
		return nil, &Err.Error{Code: Err.Econflict, Message: "Authentication required",
			Err: errors.New("user not found")}
	}
	return h.Store.User.FindByName(userName.(string))

}
