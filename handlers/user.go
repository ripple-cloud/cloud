package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"github.com/ripple-cloud/cloud/data"
)

var tokenSecret string

func init() {
	tokenSecret = os.GetEnv("TOKEN_SECRET")
	if tokenSecret == "" {
		panic("TOKEN_SECRET is not set")
	}
}

func Signup(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var respErr data.Error

		username := r.FormValue("username")
		if username == "" {
			return badRequest(w, errorMsg{"username_required", "username required"})
		}

		email := r.FormValue("email")
		if email == "" {
			return badRequest(w, errorMsg{"email_required", "email required"})
		}

		password := r.FormValue("password")
		if password == "" {
			return badRequest(w, errorMsg{"password_required", "password required"})
		}

		u := &data.User{
			Username: username,
			Email:    email,
			Password: password,
		}
		err := data.Insert(db, u)
		if err != nil {
			if err == data.ErrRecordExist {
				return badRequest(w, errorMsg{err.Code, err.Desc})
			}
			return serverError(w, r, err)
		}

		respondJSON(w, http.StatusCreated, u)
	}
}

// POST /oauth/token
// Params: grant_type, login, password
func UserToken(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		if r.FormValue("grant_type") != "password" {
			return badRequest(w, errorMsg{"unsupported_grant_type", "supports only password grant type"})
		}

		login := r.FormValue("login")
		if login == "" {
			return badRequest(w, errorMsg{"invalid_request", "login required"})
		}

		password := r.FormValue("password")
		if password == "" {
			return badRequest(w, errorMsg{"invalid_request", "password required"})
		}

		u := data.User{}
		if err := u.FindByLogin(db, login); err != nil {
			if err == data.ErrRecordNotFound {
				return badRequest(w, errorMsg{"invalid_grant", "user not found"})
			}
			return serverError(w, r, err)
		}

		if !u.Verify(db, password) {
			return badRequest(w, errorMsg{"invalid_grant", "failed to authenticate user"})
		}

		// Since all is well, generate token and add to database
		t := data.Token{
			UserID:    u.ID,
			ExpiresIn: 30 * 24 * time.Hour, // 30 days
		}
		err := data.Insert(db, t)
		if err != nil {
			return serverError(w, r, err)
		}

		// encode the token as a JSON Web token
		jt := jwt.New(jwt.SigningMethodHS256)
		jt.Claims["iat"] = t.CreatedAt.Unix()                  // issued at
		jt.Claims["exp"] = t.CreatedAt.Add(t.ExpiresIn).Unix() // expires at
		jt.Claims["jti"] = t.ID                                // token ID
		jt.Claims["user_id"] = t.UserID
		jt.Claims["scope"] = []string{"user", "hub", "app"}
		jtStr, err := jt.SignedString(tokenSecret)
		if err != nil {
			return serverError(w, r, err)
		}

		// prepare oAuth2 access token payload
		payload := struct {
			accessToken string "json:access_token"
			tokenType   string "json:token_type"
			expiresIn   string "json:expires_in"
		}{
			jtStr,
			"bearer",
			t.ExpiresIn,
		}

		respondJSON(w, http.StatusOK, payload)
	}
}
