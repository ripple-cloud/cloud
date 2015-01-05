package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	db, err := db()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := httprouter.New()

	router.POST("/signup", signupHandler(db))
	router.POST("/oauth/token", tokenHandler(db))

	router.POST("/api/v1/hub", addHubHandler(db))
	router.GET("/api/v1/hub", showHubHandler(db))
	router.DELETE("/api/v1/hub", deleteHubHandler(db))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
