package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error

	cwd, _ = os.Getwd()
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

	// Implement better routing.
	http.HandleFunc("/", signupPageHandler)
	http.HandleFunc("/loginpage", loginPageHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/home", homeHandler)

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(filepath.Join(cwd, "/github.com/ripple-cloud/cloud/templates")))))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
