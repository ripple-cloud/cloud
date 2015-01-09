package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/julienschmidt/httprouter"

	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/utils"
)

// POST /signup
// Query: username, email, password
func Signup(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var respErr data.Error

		q := map[string]string{
			"username": r.URL.Query().Get("username"),
			"email":    r.URL.Query().Get("email"),
			"password": r.URL.Query().Get("password"),
		}

		// sanitizeQuery() checks if all and only required params are included.
		respErr = utils.SanitizeQuery(r, q)
		if respErr != (data.Error{}) {
			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		user := data.User{
			Username: q["username"],
		}

		// Validate new user.
		if user.GetByUsername(db).Username == "" {
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

			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}

// POST /oauth/token
// Query: grant_type, username, password
func UserToken(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var respErr data.Error
		var user data.User

		q := map[string]string{
			"grant_type": r.URL.Query().Get("grant_type"),
			"username":   r.URL.Query().Get("username"),
			"password":   r.URL.Query().Get("password"),
		}

		user.Username = q["username"]

		// sanitizeQuery() checks if all and only required params are included.
		respErr = utils.SanitizeQuery(r, q)
		if respErr != (data.Error{}) {
			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		if q["grant_type"] != "password" {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_request",
					Description: "Invalid 'grant_type' value. 'grant_type' should be set to 'password'",
				},
			}
			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		// Check if user exists.
		if user.GetByUsername(db).Username == "" {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_client",
					Description: "invalid user credentials: user does not exist",
				},
			}

			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		} else {
			err := bcrypt.CompareHashAndPassword(user.GetByUsername(db).Password, []byte(q["password"]))
			if err != nil {
				respErr = data.Error{
					data.ErrorInfo{
						Code:        "invalid_client",
						Description: "invalid user credentials: password and user do not match",
					},
				}

				if err := utils.RespJSON(w, respErr, 400); err != nil {
					fmt.Println(err)
				}
				return
			}
		}

		// Since all is well, generate token and add to database if token has not been set.
		if user.GetByUsername(db).Token == "" {
			user.SetToken(db)
		}

		// NOTE: Left out scope and refresh_token.
		resp := data.Token{
			data.TokenInfo{
				AccessToken: user.GetByUsername(db).Token,
				TokenType:   "bearer",
				ExpiresIn:   2592000, // 30 days.
			},
		}

		if err := utils.RespJSON(w, resp, 200); err != nil {
			fmt.Println(err)
		}
	}
}
