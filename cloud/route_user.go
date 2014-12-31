package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/julienschmidt/httprouter"

	"github.com/ripple-cloud/cloud/data"
)

func signupHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		if !exist("user?", db, "username", user.Username) {
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
}

func tokenHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		if !exist("user?", db, "username", q["username"]) {
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
		if !exist("token?", db, "username", q["username"]) {
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
}
