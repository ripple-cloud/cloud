package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ripple-cloud/cloud/handlers"
	"github.com/ripple-cloud/cloud/utils"
)

func main() {
	db, err := utils.Db()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := httprouter.New()

	router.POST("/signup", handlers.Signup(db))
	router.POST("/oauth/token", handlers.UserToken(db))

	router.POST("/api/v1/hub", handlers.AddHub(db))
	router.GET("/api/v1/hub", handlers.ShowHub(db))
	router.DELETE("/api/v1/hub", handlers.DeleteHub(db))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
