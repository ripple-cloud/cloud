package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ripple-cloud/cloud/router"
)

// A Middleware handler usually changes the context and pass the request to the next handler.
// It may decide to respond early if the request can't be fulfiled (eg: authentication failure).

func SetConfig(db *sqlx.DB, tokenSecret string) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, c router.Context) error {
		c.Meta["db"] = db
		c.Meta["tokenSecret"] = tokenSecret
		return c.Next(w, r, c)
	}
}
