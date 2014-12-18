package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/ripple-cloud/cloud/data"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
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
		user.AddTo(db)
		fmt.Fprint(w, "Successful signup!")
	} else {
		fmt.Fprint(w, "Email is already taken")
	}
}

func tokenCreateHandler(w http.ResponseWriter, r *http.Request) {
	login := data.User{
		Username: r.URL.Query().Get("username"),
		Password: []byte(r.URL.Query().Get("password")),
	}

	// Validate existing user.
	if login.GetUserFrom(db).Username == login.Username {
		user := *login.GetUserFrom(db)

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password))
		if err != nil {
			fmt.Fprint(w, "Username and password do not match")
			return
		}
	} else {
		fmt.Fprint(w, "User does not exist")
		return
	}

	// After validation, add token to database.
	user := data.User{
		Username: login.Username,
		Token:    data.GenerateToken(),
	}
	user.SetToken(db)

	// Return these values in json.
	resp := data.User{
		Username: user.Username,
		Token:    user.GetUserFrom(db).Token,
	}

	w.Header().Set("Content-Type", "application/json")

	js, _ := json.Marshal(resp)
	w.Write(js)
}

func tokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	login := data.User{
		Username: r.URL.Query().Get("username"),
		Password: []byte(r.URL.Query().Get("password")),
	}

	// Validate existing user.
	if login.GetUserFrom(db).Username == login.Username {
		user := *login.GetUserFrom(db)

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password))
		if err != nil {
			fmt.Fprint(w, "Wrong password")
			return
		}
	} else {
		fmt.Fprint(w, "User does not exist")
		return
	}

	// Fetch existing token from database and return these values in json.
	resp := data.User{
		Username: login.Username,
		Token:    login.GetUserFrom(db).Token,
	}

	w.Header().Set("Content-Type", "application/json")

	js, _ := json.Marshal(resp)
	w.Write(js)
}
