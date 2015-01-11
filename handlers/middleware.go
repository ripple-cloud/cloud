package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
	res "github.com/ripple-cloud/cloud/jsonrespond"
	"github.com/ripple-cloud/cloud/router"
)

// A Middleware handler usually changes the context and pass the request to the next handler.
// It may decide to respond early if the request can't be fulfiled (eg: authentication failure).

var scopeRegex *regexp.Regexp

func init() {
	scopeRegex = regexp.MustCompile(`^(?:/api/v\d/)([^/]+)(.*)$`) // eg: /api/v0/hub/list || /api/v0/app
}

func SetDB(db *sql.DB) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, c router.Context) {
		c.Meta["db"] = db
		c.Next(w, r, c)
	}
}

func SetTokenSecret(secret string) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, c router.Context) {
		c.Meta["tokenSecret"] = secret
		c.Next(w, r, c)
	}
}

func Auth(w http.ResponseWriter, r *http.Request, c router.Context) {
	db := c.Meta["db"].(*sqlx.DB)
	tokenSecret := c.Meta["tokenSecret"].(string)

	// parse the token param
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})

	if err != nil {
		res.Unauthorized(w, res.ErrorMsg{"invalid_token", err.Error()})
		return
	}

	// check if the token is eligible for current scope
	scope := scopeRegex.FindStringSubmatch(r.URL.Path)[1]
	scopes := token.Claims["scopes"].([]string)

	if !contains(scopes, scope) {
		res.Forbidden(w, res.ErrorMsg{"invalid_scope", "token is not valid for scope"})
	}

	// check if the token was revoked from DB
	t := data.Token{}
	err = t.Get(db, token.Claims["jti"].(int64))
	if err != nil {
		if _, ok := err.(*data.Error); ok {
			res.Unauthorized(w, res.ErrorMsg{"invalid_token", "token is not valid"})
			return
		}
		log.Print("[error] auth: %s", err)
		res.ServerError(w, res.ErrorMsg{"internal_server_error", "Something went wrong"})
		return
	}
	if t.RevokedAt != nil {
		res.Unauthorized(w, res.ErrorMsg{"invalid_token", "token is not valid"})
		return
	}

	// valid token
	// set the user id to context and pass to next handler
	c.Meta["user_id"] = t.UserID

	c.Next(w, r, c)
}

func contains(col []string, val string) bool {
	for _, cur := range col {
		if cur == val {
			return true
		}
	}
	return false
}
