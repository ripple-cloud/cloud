package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ripple-cloud/cloud/router"
)

// A Middleware handler usually changes the context and pass the request to the next handler.
// It may decide to respond early if the request can't be fulfiled (eg: authentication failure).

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
