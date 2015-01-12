package jsonrespond

import (
	"encoding/json"
	"net/http"
)

type ErrorMsg struct {
	Error     string "json:error"
	ErrorDesc string "json:error_description"
}

func Respond(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)

	o, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = w.Write(o)
	return err
}

func OK(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusOK, payload)
}

func BadRequest(w http.ResponseWriter, err ErrorMsg) error {
	return Respond(w, http.StatusBadRequest, err)
}

func Unauthorized(w http.ResponseWriter, err ErrorMsg) error {
	return Respond(w, http.StatusUnauthorized, err)
}

func Forbidden(w http.ResponseWriter, err ErrorMsg) error {
	return Respond(w, http.StatusForbidden, err)
}

func ServerError(w http.ResponseWriter, err ErrorMsg) error {
	return Respond(w, http.StatusInternalServerError, err)
}
