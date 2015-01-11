package handlers

//
// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
//
// 	"github.com/julienschmidt/httprouter"
// 	"github.com/ripple-cloud/cloud/data"
// 	"github.com/ripple-cloud/cloud/utils"
// )
//
// // POST /api/add
// // Query: hub, token.
// func AddHub(db *sql.DB) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		var respErr data.Error
// 		var user data.User
//
// 		q := map[string]string{
// 			"hub":   r.URL.Query().Get("hub"),
// 			"token": r.URL.Query().Get("token"),
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
// 		hub := data.Hub{
// 			Hub:    q["hub"],
// 			UserID: user.GetByToken(db).ID,
// 		}
//
// 		//check
// 		if hub.GetByHub(db).Hub == "" {
// 			respErr = data.Error{
// 				data.ErrorInfo{
// 					Code:        "invalid_request",
// 					Description: "hub already exists",
// 				},
// 			}
//
// 			if err := utils.RespJSON(w, respErr, 400); err != nil {
// 				fmt.Println(err)
// 			}
// 			return
// 		}
// 		// FIXME: error message --> pq: duplicate key value violates unique constraint "hubs_hub_key"
// 		hub.Add(db)
//
// 		resp := data.AddHub{
// 			data.AddHubInfo{
// 				ID:     hub.GetByHub(db)[0].ID,
// 				UserID: hub.UserID,
// 				Slug:   q["hub"],
// 			},
// 		}
//
// 		if err := utils.RespJSON(w, resp, 201); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
//
// // GET /api/hub
// // Query: "token"
// func ShowHub(db *sql.DB) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		var respErr data.Error
// 		var user data.User
// 		var hub data.Hub
//
// 		q := map[string]string{
// 			"token": r.URL.Query().Get("token"),
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
// 		resp := data.ShowHub{
// 			Hubs: hub.Get(db, "user_id", user.GetByToken(db).ID),
// 		}
//
// 		if err := utils.RespJSON(w, resp, 200); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
//
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
