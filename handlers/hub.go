package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/utils"
)

func AddHub(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// POST /api/add
		// Query: hub, token.
		var respErr data.Error
		var user data.User

		q := map[string]string{
			"hub":   r.URL.Query().Get("hub"),
			"token": r.URL.Query().Get("token"),
		}

		respErr = utils.SanitizeQuery("add", r, q)
		if respErr != (data.Error{}) {
			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		if !utils.Exist("token?", db, "token", q["token"]) {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_client",
					Description: "token does not exist",
				},
			}

			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		hub := data.Hub{
			Hub:    q["hub"],
			UserID: user.Get(db, "token", q["token"]).ID,
		}

		if utils.Exist("hub?", db, q["hub"], hub.UserID) {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_request",
					Description: "hub already exists",
				},
			}

			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}
		// FIXME: error message --> pq: duplicate key value violates unique constraint "hubs_hub_key"
		hub.Add(db)

		resp := data.AddHub{
			data.AddHubInfo{
				ID:     hub.Get(db, "hub", q["hub"])[0].ID,
				UserID: hub.UserID,
				Slug:   q["hub"],
			},
		}

		if err := utils.RespJSON(w, resp, 201); err != nil {
			fmt.Println(err)
		}
	}
}

func ShowHub(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// GET /api/hub
		// Query: "token"
		var respErr data.Error
		var user data.User
		var hub data.Hub

		q := map[string]string{
			"token": r.URL.Query().Get("token"),
		}

		respErr = utils.SanitizeQuery("get", r, q)
		if respErr != (data.Error{}) {
			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		if !utils.Exist("token?", db, "token", q["token"]) {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_client",
					Description: "token does not exist",
				},
			}

			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		resp := data.ShowHub{
			Hubs: hub.Get(db, "user_id", user.Get(db, "token", q["token"]).ID),
		}

		if err := utils.RespJSON(w, resp, 200); err != nil {
			fmt.Println(err)
		}
	}
}

func DeleteHub(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// DELETE /api/hub
		// Query: token, id
		var respErr data.Error
		var user data.User
		var hub data.Hub

		q := map[string]string{
			"token": r.URL.Query().Get("token"),
			"id":    r.URL.Query().Get("id"),
		}

		respErr = utils.SanitizeQuery("delete", r, q)
		if respErr != (data.Error{}) {
			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		if !utils.Exist("token?", db, "token", q["token"]) {
			respErr = data.Error{
				data.ErrorInfo{
					Code:        "invalid_client",
					Description: "token does not exist",
				},
			}

			if err := utils.RespJSON(w, respErr, 400); err != nil {
				fmt.Println(err)
			}
			return
		}

		// Check that hub id belongs to user of access token.
		userIDToken := user.Get(db, "token", q["token"]).ID
		userIDID := hub.Get(db, "id", q["id"])[0].UserID

		if userIDToken != userIDID {
			resp := data.Error{
				data.ErrorInfo{
					Code:        "Invalid_request",
					Description: "hub does not belong to user",
				},
			}
			if err := utils.RespJSON(w, resp, 400); err != nil {
				fmt.Println(err)
			}
		}

		resp := data.DeleteHub{
			data.DeleteHubInfo{
				ID:     q["id"],
				Hub:    hub.Get(db, "id", q["id"])[0].Hub,
				UserID: userIDID,
			},
		}

		hub.Delete(db, "id", q["id"])
		if err := utils.RespJSON(w, resp, 200); err != nil {
			fmt.Println(err)
		}
	}
}