package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ripple-cloud/cloud/data"
)

//
// Helper functions for cloud.go, route_user.go, route_hub.go.
// respJSON(), exist(), sanitizeQuery().
//

func RespJSON(w http.ResponseWriter, resp interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)

	js, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(js)

	return nil
}

func Exist(obj string, db *sql.DB, col, value string) bool {
	user := data.User{}
	hub := data.Hub{}

	switch obj {
	case "user?":
		if user.Get(db, col, value).Username == "" {
			return false
		}
	case "token?":
		if user.Get(db, col, value).Token == "" {
			return false
		}
	case "hub?":
		for i := 0; i < len(hub.Get(db, "user_id", value)); i++ {
			if hub.Get(db, "user_id", value)[i].Hub == col {
				return true
			} else {
				return false
			}
		}
	default:
		fmt.Println("Error: check that you're querying the right objects")
	}

	return true
}

// FIXME: This should not require an action argument. (let the caller handle the error)
func SanitizeQuery(r *http.Request, q map[string]string) data.Error {
	var respErr data.Error

	for k, v := range q {
		if v == "" {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_request",
					Description: fmt.Sprintf("Missing parameter, %s is required", k),
				},
			}
			return respErr
		}
	}

	a := 0
	for key := range r.URL.Query() {
		for k, _ := range q {
			if key == k {
				a += 1
			}
		}
		if a == 0 {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_request",
					Description: fmt.Sprintf("Invalid parameter %s", key),
				},
			}
			return respErr
		}
		a = 0
	}

	return respErr
}
