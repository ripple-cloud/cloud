package handlers

import (
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
)

// POST api/v0/app/:slug
// Params: hub_id
func RegisterApp() error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}

	slug := c.Params.ByName("slug")
	if slug == "" {
		return res.BadRequest(w, res.ErrorMsg{"slug_required", "slug required"})
	}

	// TODO: In future we can support registering an app on all hubs with a singe request
	hubID := r.FormValue("hub_id")
	if hubID == "" {
		return res.BadRequest(w, res.ErrorMsg{"hub_id_required", "hub id required"})
	}

	a := &data.App{
		Slug:  slug,
		HubID: hub_id,
	}
	if err := a.Insert(db); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.BadRequest(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	return res.Created(w, a)
}

// DELETE api/v0/app/:slug
// Params: hub_id
func DeregisterApp() error {
	db, ok := c.Meta["db"].(*sqlx.DB)
	if !ok {
		return errors.New("db not set in context")
	}

	slug := c.Params.ByName("slug")
	if slug == "" {
		return res.BadRequest(w, res.ErrorMsg{"slug_required", "slug required"})
	}

	// TODO: In future we can support de-registering an app from all hubs with a singe request
	hubID := r.FormValue("hub_id")
	if hubID == "" {
		return res.BadRequest(w, res.ErrorMsg{"hub_id_required", "hub id required"})
	}

	a := &data.App{}
	if err := a.Get(db, slug, hubID); err != nil {
		if e, ok := err.(*data.Error); ok {
			return res.NotFound(w, res.ErrorMsg{e.Code, e.Desc})
		}
		return err
	}

	return res.Respond(w, http.StatusOK, a)
}
