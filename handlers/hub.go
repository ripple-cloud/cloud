package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ripple-cloud/cloud/data"
	res "github.com/ripple-cloud/cloud/jsonrespond"
	"github.com/ripple-cloud/cloud/router"
)

// POST /api/v0/hub
// Params: access_token, slug, (scope?)
func AddHub(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, _ := c.Meta["db"].(*sqlx.DB)

	slug := r.FormValue("slug")
	if slug == "" {
		return res.BadRequest(w, res.ErrorMsg{"invalid_request", "slug required"})
	}

	// Since all is well, add hub to database
	h := data.Hub{
		Slug:   slug,
		UserID: c.Meta["user_id"].(int64),
	}
	if err := h.Insert(db); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	return res.OK(w, h)
}

// GET /api/v0/hub
// Params: access_token
func ShowHub(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, _ := c.Meta["db"].(*sqlx.DB)

	// Since all is well, get hub(s) from database
	var h data.Hubs
	if err := h.SelectByUserId(db, c.Meta["user_id"].(int64)); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	// prepare oAuth2 access token payload
	payload := struct {
		Hubs []string `json:"hub"`
	}{
		h,
	}

	return res.OK(w, payload)
}

// DELETE /api/v0/hub
// Params: access_token, slug
func DeleteHub(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, _ := c.Meta["db"].(*sqlx.DB)
	userid := c.Meta["user_id"].(int64)

	slug := r.FormValue("slug")
	if slug == "" {
		return res.BadRequest(w, res.ErrorMsg{"invalid_request", "slug required"})
	}
	h := data.Hub{}
	if err := h.Get(db, slug); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	if userid != h.UserID {
		return res.BadRequest(w, res.ErrorMsg{"invalid_request", "user does not own hub"})
	}

	// Since all is well, delete hub from database
	h = data.Hub{
		Slug: slug,
	}
	if err := h.Delete(db); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	return res.OK(w, h)
}
