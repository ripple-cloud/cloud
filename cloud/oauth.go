package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/julienschmidt/httprouter"

	"github.com/ripple-cloud/cloud/data"
)

//
// Handlers
//

func signupHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// POST /signup
	// Query: username, email, password
	var respErr data.Error

	q := map[string]string{
		"username": r.URL.Query().Get("username"),
		"email":    r.URL.Query().Get("email"),
		"password": r.URL.Query().Get("password"),
	}

	// sanitizeQuery() checks if all and only required params are included.
	respErr = sanitizeQuery("signup", r, q)
	if respErr != (data.Error{}) {
		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	user := data.User{
		Username: q["username"],
	}

	// Validate new user.
	if !exist("user?", "username", user.Username) {
		user := data.User{
			Username:  user.Username,
			Email:     q["email"],
			Password:  data.Encrypt(q["password"]),
			Token:     "",
			CreatedAt: time.Now(),
		}
		user.Add(db)
		// TODO: render JSON
		fmt.Fprint(w, "Successful signup!")

	} else {
		respErr = data.Error{
			data.ErrorInfo{
				Code:        "invalid_client",
				Description: "username is already taken",
			},
		}
		//TODO: Check if email is unique and add error handling.

		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// POST api/oauth/token
	// Query: grant_type, username, password
	var respErr data.Error
	var user data.User

	q := map[string]string{
		"grant_type": r.URL.Query().Get("grant_type"),
		"username":   r.URL.Query().Get("username"),
		"password":   r.URL.Query().Get("password"),
	}

	// sanitizeQuery() checks if (i) all and only required params are included (ii) grant_type is set to password.
	respErr = sanitizeQuery("token", r, q)
	if respErr != (data.Error{}) {
		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	// Check if user exists.
	if !exist("user?", "username", q["username"]) {
		respErr = data.Error{
			data.ErrorInfo{
				Code:        "invalid_client",
				Description: "invalid user credentials: user does not exist",
			},
		}

		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	} else {
		err := bcrypt.CompareHashAndPassword(user.Get(db, "username", q["username"]).Password, []byte(q["password"]))
		if err != nil {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_client",
					Description: "invalid user credentials: password and user do not match",
				},
			}

			if err := respJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}
	}

	// Since all is well, generate token and add to database if token has not been set.
	if !exist("token?", "username", q["username"]) {
		user.SetToken(db, "username", q["username"])
	}

	// NOTE: Left out scope and refresh_token.
	resp := data.Token{
		data.TokenInfo{
			AccessToken: user.Get(db, "username", q["username"]).Token,
			TokenType:   "bearer",
			ExpiresIn:   2592000, // 30 days.
		},
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := respJSON(w, resp, 200); err != nil {
		fmt.Println(err)
	}
}

func addHubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// POST /api/add
	// Query: hub, token.
	var respErr data.Error
	var user data.User

	q := map[string]string{
		"hub":   r.URL.Query().Get("hub"),
		"token": r.URL.Query().Get("token"),
	}

	respErr = sanitizeQuery("add", r, q)
	if respErr != (data.Error{}) {
		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	if !exist("token?", "token", q["token"]) {
		respErr = data.Error{
			data.ErrorInfo{
				Code:        "invalid_client",
				Description: "token does not exist",
			},
		}

		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	hub := data.Hub{
		Hub:    q["hub"],
		UserID: user.Get(db, "token", q["token"]).ID,
	}

	if exist("hub?", q["hub"], hub.UserID) {
		respErr = data.Error{
			data.ErrorInfo{
				Code:        "invalid_request",
				Description: "hub already exists",
			},
		}

		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}
	// FIXME: error message --> pq: duplicate key value violates unique constraint "hubs_hub_key"
	hub.Add(db)

	resp := data.AddHub{
		data.AddHubInfo{
			ID:     hub.Get(db, "hub", q["hub"])[0].ID,
			UserID: hub.UserID,
			Slug:   q["hub"],
		},
	}

	if err := respJSON(w, resp, 201); err != nil {
		fmt.Println(err)
	}
}

func showHubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// GET /api/hub
	// Query: "token"
	var respErr data.Error
	var user data.User
	var hub data.Hub

	q := map[string]string{
		"token": r.URL.Query().Get("token"),
	}

	respErr = sanitizeQuery("get", r, q)
	if respErr != (data.Error{}) {
		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	if !exist("token?", "token", q["token"]) {
		respErr = data.Error{
			data.ErrorInfo{
				Code:        "invalid_client",
				Description: "token does not exist",
			},
		}

		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	resp := data.ShowHub{
		Hubs: hub.Get(db, "user_id", user.Get(db, "token", q["token"]).ID),
	}

	if err := respJSON(w, resp, 200); err != nil {
		fmt.Println(err)
	}
}

func deleteHubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// DELETE /api/hub
	// Query: token, id
	var respErr data.Error
	var user data.User
	var hub data.Hub

	q := map[string]string{
		"token": r.URL.Query().Get("token"),
		"id":    r.URL.Query().Get("id"),
	}

	respErr = sanitizeQuery("delete", r, q)
	if respErr != (data.Error{}) {
		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	if !exist("token?", "token", q["token"]) {
		respErr = data.Error{
			data.ErrorInfo{
				Code:        "invalid_client",
				Description: "token does not exist",
			},
		}

		if err := respJSON(w, respErr, 400); err != nil {
			fmt.Println(err)
		}
		return
	}

	// Check that hub id belongs to user of access token.
	userIDToken := user.Get(db, "token", q["token"]).ID
	userIDID := hub.Get(db, "id", q["id"])[0].UserID

	if userIDToken != userIDID {
		resp := data.Error{
			data.ErrorInfo{
				Code:        "Invalid_request",
				Description: "hub does not belong to user",
			},
		}
		if err := respJSON(w, resp, 400); err != nil {
			fmt.Println(err)
		}
	}

	resp := data.DeleteHub{
		data.DeleteHubInfo{
			ID:     q["id"],
			Hub:    hub.Get(db, "id", q["id"])[0].Hub,
			UserID: userIDID,
		},
	}

	hub.Delete(db, "id", q["id"])
	if err := respJSON(w, resp, 200); err != nil {
		fmt.Println(err)
	}
}

//
// Helper functions.
// respJSON(), exist(), sanitizeQuery().
//

func respJSON(w http.ResponseWriter, resp interface{}, code int) error {
	w.WriteHeader(code)

	js, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(js)

	return nil
}

func exist(obj, col, value string) bool {
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

func sanitizeQuery(action string, r *http.Request, q map[string]string) data.Error {
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

	switch action {
	case "token":
		if q["grant_type"] != "password" {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_request",
					Description: "Invalid 'grant_type' value. 'grant_type' should be set to 'password'",
				},
			}
			return respErr
		}
	}
	return respErr
}
