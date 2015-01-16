package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ripple-cloud/cloud/data"
	res "github.com/ripple-cloud/cloud/jsonrespond"
	"github.com/ripple-cloud/cloud/router"
)

// POST /signup
// Params: username, email, password
func Signup(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}

	username := r.FormValue("username")
	if username == "" {
		return res.BadRequest(w, res.ErrorMsg{"username_required", "username required"})
	}

	email := r.FormValue("email")
	if email == "" {
		return res.BadRequest(w, res.ErrorMsg{"email_required", "email required"})
	}

	password := r.FormValue("password")
	if password == "" {
		return res.BadRequest(w, res.ErrorMsg{"password_required", "password required"})
	}

	u := &data.User{
		Username: username,
		Email:    email,
	}
	if err := u.EncryptPassword(password); err != nil {
		return err
	}
	if err := u.Insert(db); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	return res.Respond(w, http.StatusCreated, u)
}

// POST /oauth/token
// Params: grant_type, login, password
// Requires a tokenSecret to be set in context
func UserToken(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}
	tokenSecret, ok := c.Meta["tokenSecret"].([]byte)
	if !ok {
		return errors.New("token secret not set in context")
	}

	if r.FormValue("grant_type") != "password" {
		return res.BadRequest(w, res.ErrorMsg{"unsupported_grant_type", "supports only password grant type"})
	}

	login := r.FormValue("login")
	if login == "" {
		return res.BadRequest(w, res.ErrorMsg{"invalid_request", "login required"})
	}

	password := r.FormValue("password")
	if password == "" {
		return res.BadRequest(w, res.ErrorMsg{"invalid_request", "password required"})
	}

	u := data.User{}
	if err := u.GetByLogin(db, login); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{"invalid_grant", e.Desc})
		}
		return err
	}

	if !u.VerifyPassword(password) {
		return res.BadRequest(w, res.ErrorMsg{"invalid_grant", "failed to authenticate user"})
	}

	// Since all is well, generate token and add to database
	t := data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(), // 30 days
	}
	if err := t.Insert(db); err != nil {
		return err
	}

	// get the encoded JSON Web token
	jwt, err := t.EncodeJWT(tokenSecret)
	if err != nil {
		return err
	}

	// prepare oAuth2 access token payload
	payload := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   string `json:"expires_in"`
	}{
		jwt,
		"bearer",
		time.Duration(t.ExpiresIn).String(),
	}

	return res.OK(w, payload)
}
