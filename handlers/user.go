package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ripple-cloud/cloud/data"
	res "github.com/ripple-cloud/cloud/jsonrespond"
	"github.com/ripple-cloud/cloud/router"
)

var tokenSecret string

func init() {
	tokenSecret = os.GetEnv("TOKEN_SECRET")
	if tokenSecret == "" {
		panic("TOKEN_SECRET is not set")
	}
}

// POST /signup
// Params: username, email, password
func Signup(w http.ResponseWriter, r *http.Request, c router.Context) {
	db, ok := c.Meta["db"].(*sql.DB)
	if !ok {
		log.Print("[error] signup: DB not set in context")
		res.ServerError(w, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}

	username := r.FormValue("username")
	if username == "" {
		res.BadRequest(w, errorMsg{"username_required", "username required"})
		return
	}

	email := r.FormValue("email")
	if email == "" {
		res.BadRequest(w, errorMsg{"email_required", "email required"})
		return
	}

	password := r.FormValue("password")
	if password == "" {
		res.BadRequest(w, errorMsg{"password_required", "password required"})
		return
	}

	u := &data.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	err := data.Insert(db, u)
	if err != nil {
		if err == data.ErrRecordExist {
			res.BadRequest(w, errorMsg{err.Code, err.Desc})
			return
		}
		log.Printf("[error] signup: %s", err)
		res.ServerError(w, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}

	res.Respond(w, http.StatusCreated, u)
}

// POST /oauth/token
// Params: grant_type, login, password
// Requires a tokenSecret to be set in context
func UserToken(w http.ResponseWriter, r *http.Request, c router.Context) {
	db, ok := c.Meta["db"].(*sql.DB)
	if !ok {
		log.Print("[error] signup: DB not set in context")
		res.ServerError(w, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}
	tokenSecret, ok := c.Meta["tokenSecret"].(string)
	if !ok {
		log.Print("[error] userToken: token secret not set in context")
		res.ServerError(w, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}

	if r.FormValue("grant_type") != "password" {
		res.BadRequest(w, errorMsg{"unsupported_grant_type", "supports only password grant type"})
		return
	}

	login := r.FormValue("login")
	if login == "" {
		res.BadRequest(w, errorMsg{"invalid_request", "login required"})
		return
	}

	password := r.FormValue("password")
	if password == "" {
		res.BadRequest(w, errorMsg{"invalid_request", "password required"})
		return
	}

	u := data.User{}
	if err := u.FindByLogin(db, login); err != nil {
		if err == data.ErrRecordNotFound {
			res.BadRequest(w, errorMsg{"invalid_grant", "user not found"})
			return
		}
		log.Print("[error] userToken: %s", err)
		res.ServerError(w, r, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}

	if !u.Verify(db, password) {
		return res.BadRequest(w, errorMsg{"invalid_grant", "failed to authenticate user"})
	}

	// Since all is well, generate token and add to database
	t := data.Token{
		UserID:    u.ID,
		ExpiresIn: 30 * 24 * time.Hour, // 30 days
	}
	err := data.Insert(db, t)
	if err != nil {
		log.Print("[error] userToken: %s", err)
		res.ServerError(w, r, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}

	// encode the token as a JSON Web token
	jt := jwt.New(jwt.SigningMethodHS256)
	jt.Claims["iat"] = t.CreatedAt.Unix()                  // issued at
	jt.Claims["exp"] = t.CreatedAt.Add(t.ExpiresIn).Unix() // expires at
	jt.Claims["jti"] = t.ID                                // token ID
	jt.Claims["user_id"] = t.UserID
	jt.Claims["scopes"] = []string{"user", "hub", "app"}
	jtStr, err := jt.SignedString(tokenSecret)
	if err != nil {
		log.Print("[error] userToken: %s", err)
		res.ServerError(w, r, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
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

	res.Respond(w, http.StatusOK, payload)
}
