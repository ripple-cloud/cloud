package handlers

import (
	"encoding/json"
	"net/http"
)

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
