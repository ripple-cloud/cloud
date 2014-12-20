package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	r.HandleFunc("/signup", signupHandler).Methods("POST")
	r.HandleFunc("/api/oauth/token", tokenHandler).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
