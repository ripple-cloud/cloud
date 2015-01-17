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
		return err
	}

	// prepare oAuth2 access token payload
	payload := struct {
		Slug string `json:"slug"`
	}{
		slug,
	}

	return res.OK(w, payload)
}

// GET /api/v0/hub
// Params: access_token
func ShowHub(w http.ResponseWriter, r *http.Request, c router.Context) error {
	db, _ := c.Meta["db"].(*sqlx.DB)

	// Since all is well, get hub(s) from database
	var h data.Hubs
	if err := h.GetByUserId(db, c.Meta["user_id"].(int64)); err != nil {
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

// // DELETE /api/hub
// // Query: token, id
// func DeleteHub(db *sql.DB) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		var respErr data.Error
// 		var user data.User
// 		var hub data.Hub
//
// 		q := map[string]string{
// 			"token": r.URL.Query().Get("token"),
// 			"id":    r.URL.Query().Get("id"),
// 		}
//
// 		respErr = utils.SanitizeQuery(r, q)
// 		if respErr != (data.Error{}) {
// 			if err := utils.RespJSON(w, respErr, 400); err != nil {
// 				fmt.Println(err)
// 			}
// 			return
// 		}
//
// 		if user.GetByToken(db).Token == "" {
// 			respErr = data.Error{
// 				data.ErrorInfo{
// 					Code:        "invalid_client",
// 					Description: "token does not exist",
// 				},
// 			}
//
// 			if err := utils.RespJSON(w, respErr, 400); err != nil {
// 				fmt.Println(err)
// 			}
// 			return
// 		}
//
// 		// Check that hub id belongs to user of access token.
// 		userIDToken := user.GetByToken(db).ID
// 		//check
// 		userIDID := hub.GetByID(db).UserID
//
// 		if userIDToken != userIDID {
// 			resp := data.Error{
// 				data.ErrorInfo{
// 					Code:        "Invalid_request",
// 					Description: "hub does not belong to user",
// 				},
// 			}
// 			if err := utils.RespJSON(w, resp, 400); err != nil {
// 				fmt.Println(err)
// 			}
// 		}
//
// 		resp := data.DeleteHub{
// 			data.DeleteHubInfo{
// 				ID:     q["id"],
// 				Hub:    hub.GetByID(db).Hub,
// 				UserID: userIDID,
// 			},
// 		}
//
// 		//define hub struct id
//
// 		hub.Delete(db)
// 		if err := utils.RespJSON(w, resp, 200); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
