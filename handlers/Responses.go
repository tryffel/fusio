package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ResponseBody map[string]interface{}

const (
	ResponseAccessDenied   = "Access denied"
	ResponseInvalidJson    = "Invalid json"
	ResponseInvalidBody    = "Invalid request body"
	ResponseInternalError  = "Internal error"
	ResponsePasswordFailed = "Invalid username or password"
	ResponseUnauthorized   = "Forbidden"
	ResponseInvalidId      = "Invalid ID"
	ResponseResourceExists = "Resourse exists"
	ResponseNotFound       = "Not found"
	ResponseInvalidToken   = "Invalid token"
	ResponseStatus         = "status"
	ResponseOk             = "ok"
	ResponseCreated        = "created"
	ResponseUpdated        = "updated"
	ResponseDeleted        = "deleted"
)

func JsonErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "Application/json")
	http.Error(w, fmt.Sprintf("{\"Error\": \"%s\"}", message), code)

}

func InvalidJsonResponse(w http.ResponseWriter) {
	JsonErrorResponse(w, "Invalid json", http.StatusBadRequest)
}

// JsonResponse Write json response with statuscode
func JsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "Application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logrus.Error(ResponseInvalidJson, err)
	}
}

// JsonMessage Return single message field with status code
func JsonMessage(w http.ResponseWriter, key string, value string) {
	data := make(map[string]string)
	data[key] = value
	JsonResponse(w, data)
}

func JsonResponseOk(w http.ResponseWriter) {
	JsonMessage(w, ResponseStatus, ResponseOk)
	return
}

func JsonResponseCreated(w http.ResponseWriter, data *ResponseBody) {
	if data == nil {
		JsonMessage(w, ResponseStatus, ResponseCreated)
	} else {
		(*data)[ResponseStatus] = ResponseCreated
		JsonResponse(w, data)
	}
}

func JsonResponseUpdated(w http.ResponseWriter, data *ResponseBody) {
	if data == nil {
		JsonMessage(w, ResponseStatus, ResponseUpdated)
	} else {
		(*data)[ResponseStatus] = ResponseUpdated
		JsonResponse(w, data)
	}
}

func JsonResponseDeleted(w http.ResponseWriter) {
	JsonMessage(w, ResponseStatus, ResponseDeleted)
}

func DtoFromRequest(r *http.Request, value interface{}) error {
	return json.NewDecoder(r.Body).Decode(value)
}
