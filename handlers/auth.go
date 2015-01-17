package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
	res "github.com/ripple-cloud/cloud/jsonrespond"
	"github.com/ripple-cloud/cloud/router"
)

var scopeRegex *regexp.Regexp

func init() {
	scopeRegex = regexp.MustCompile(`^(?:/api/v\d/)([^/]+)(.*)$`) // eg: /api/v0/hub/list || /api/v0/app
}

func Auth(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}
	tokenSecret, ok := c.Meta["tokenSecret"].([]byte)
	if !ok {
		return errors.New("token secret not set in context")
	}

	// parse the token param
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil {
		return res.Unauthorized(w, res.ErrorMsg{"invalid_token", err.Error()})
	}

	// check if the token is eligible for current scope
	scope := scopeRegex.FindStringSubmatch(r.URL.Path)[1]
	scopes := token.Claims["scopes"].([]interface{})

	if !contains(scopes, scope) {
		return res.Forbidden(w, res.ErrorMsg{"invalid_scope", "token is not valid for this scope"})
	}

	// check if the token was revoked from DB
	t := data.Token{}
	err = t.Get(db, int64(token.Claims["jti"].(float64)))
	if err != nil {
		if _, ok := err.(*data.Error); ok {
			return res.Unauthorized(w, res.ErrorMsg{"invalid_token", "token is not valid"})
		}
		return err
	}
	if t.RevokedAt != nil {
		return res.Unauthorized(w, res.ErrorMsg{"invalid_token", "token is not valid"})
	}

	// valid token
	// set the user id to context and pass to next handler
	c.Meta["user_id"] = t.UserID

	return c.Next(w, r, c)
}

func contains(col []interface{}, val string) bool {
	for _, cur := range col {
		if cur == val {
			return true
		}
	}
	return false
}
