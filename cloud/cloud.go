package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func main() {
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := httprouter.New()

	r.POST("/signup", signupHandler)
	r.POST("/api/oauth/token", tokenHandler)

	r.POST("/api/hub", addHubHandler)
	r.GET("/api/hub", showHubHandler)
	r.DELETE("/api/hub", deleteHubHandler)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
