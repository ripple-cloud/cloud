package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ripple-cloud/cloud/handlers"
	"github.com/ripple-cloud/cloud/router"
)

var dbURL, tokenSecret string

func init() {
	dbURL = os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB_URL not set")
	}

	tokenSecret = os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		panic("TOKEN_SECRET is not set")
	}
}

func main() {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.New()

	// default handlers are applied to all routes
	r.Default(handlers.SetConfig(db, tokenSecret))

	// unauthenticated routes
	r.POST("/signup", handlers.Signup)
	r.POST("/oauth/token", handlers.UserToken)

	// authenticated routes
	// r.POST("/api/v0/hub", handlers.Auth, handlers.AddHub)
	// r.GET("/api/v0/hub", handlers.Auth, handlers.ShowHub)
	// r.DELETE("/api/v0/hub", handlers.Auth, handlers.DeleteHub)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
