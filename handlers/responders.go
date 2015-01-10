package handlers

import (
	"encoding/json"
	"net/http"
)

type errorMsg struct {
	error     string "json:error"
	errorDesc string "json:error_description"
}

func badRequest(w http.ResponseWriter, err errorMsg) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)

	o, err := json.Marshal(err)
	if err != nil {
		log.Infof("[error] badRequest: %s", err)
		return
	}
	w.Wrtite(o)
}

// TODO: server errors must be sent to something like AirBrake (in production)
func serverError(w http.ResponseWriter, r *http.Request, err error) {
	// log the error for internal reference
	log.Infof("[error] Server Error (%s): %s ", r.RequestURI, err)

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	o, err := json.Marshal(errorMsg{"internal_server_error", "Something went wrong"})
	if err != nil {
		log.Infof("[error] serverError: %s", err)
		return
	}
	w.Wrtite(o)
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)

	o, err := json.Marshal(payload)
	if err != nil {
		log.Infof("[error] respondJSON: %s", err)
		return
	}
	w.Wrtite(o)
}
