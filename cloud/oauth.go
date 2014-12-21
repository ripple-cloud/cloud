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

func signupHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	login := data.User{
		Username: r.URL.Query().Get("username"),
	}

	// Validate new user.
	if login.GetUserFrom(db).Username == "" {
		user := data.User{
			Username:  login.Username,
			Email:     r.URL.Query().Get("email"),
			Password:  data.Encrypt(r.URL.Query().Get("password")),
			Token:     "",
			CreatedAt: time.Now(),
		}
		err := user.AddTo(db)
		if err != nil {
			return
		}
		fmt.Fprint(w, "Successful signup!")

	} else {
		fmt.Fprint(w, "Email is already taken")
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var respError data.ErrorResponse

	// Sanitize query params.
	for k := range r.URL.Query() {
		if k != "grant_type" && k != "username" && k != "password" {
			w.WriteHeader(400)

			respError.Error = "Invalid_request"
			respError.Description = fmt.Sprintf("Invalid parameter %s", k)

			js, _ := json.Marshal(respError)
			w.Write(js)
			return
		}
	}

	grant_type := r.URL.Query().Get("grant_type")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// Ensure required params are included in request.
	if grant_type == "" || username == "" || password == "" {
		w.WriteHeader(400)
		respError.Error = "Invalid_request"

		switch {
		case grant_type == "":
			respError.Description = "Missing parameter. 'grant_type' is required"
		case username == "":
			respError.Description = "Missing parameter. 'username' is required"
		case password == "":
			respError.Description = "Missing parameter. 'password' is required"
		}

		js, _ := json.Marshal(respError)
		w.Write(js)
		return
	}

	// grant_type can only be set to password.
	if grant_type != "password" {
		w.WriteHeader(400)

		respError.Error = "Invalid_request"
		respError.Description = "Invalid 'grant_type' value. 'grant_type' should be set to 'password'"

		js, _ := json.Marshal(respError)
		w.Write(js)
		return
	}

	// Execute the rest of code if query params are valid.
	login := data.User{
		Username: username,
		Password: []byte(password),
	}

	// Validate existing user.
	if login.GetUserFrom(db).Username == login.Username {
		user := *login.GetUserFrom(db)

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password))
		if err != nil {
			w.WriteHeader(400)

			respError.Error = "Invalid_client"
			respError.Description = "Invalid user credentials: password does not match user"

			js, _ := json.Marshal(respError)
			w.Write(js)
			return
		}
	} else {
		w.WriteHeader(400)

		respError.Error = "Invalid_client"
		respError.Description = "Invalid user credentials: User does not exist"

		js, _ := json.Marshal(respError)
		w.Write(js)
		return
	}

	// After validation, generate token and add to database if token has not been set.
	user := data.User{
		Username: login.Username,
		Token:    data.GenerateToken(),
	}
	if user.GetUserFrom(db).Token == "" {
		user.SetToken(db)
	}

	// Return token in json.
	// NOTE: Left out scope and refresh_token.
	resp := data.TokenResponse{
		AccessToken: user.GetUserFrom(db).Token,
		TokenType:   "bearer",
		ExpiresIn:   2592000, // 30 days.
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	js, _ := json.Marshal(resp)
	w.Write(js)
}
