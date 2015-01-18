package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ripple-cloud/cloud/dispatcher"
	"github.com/ripple-cloud/cloud/handlers"
	"github.com/ripple-cloud/cloud/router"
)

var dbURL, tokenSecret, addr string

func init() {
	dbURL = os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB_URL not set")
	}

	tokenSecret = os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		panic("TOKEN_SECRET is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // defaults to port 3000
	}
	addr = "0.0.0.0:" + port
}

func main() {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// connect to MQTT broker
	b := broker.NewMQTTBroker()
	if err := b.Connect(broker); err != nil {
		log.Fatal(err)
	}
	defer b.Disconnect()

	err := dispatcher.Start(db, broker)

	r := router.New()
	// default handlers are applied to all routes
	r.Default(handlers.SetConfig(db, []byte(tokenSecret)))

	// unauthenticated routes
	r.POST("/signup", handlers.Signup)
	r.POST("/oauth/token", handlers.UserToken)

	// authenticated routes
	r.POST("/api/v0/hub", handlers.Auth, handlers.AddHub)
	r.GET("/api/v0/hub", handlers.Auth, handlers.ShowHub)
	r.DELETE("/api/v0/hub", handlers.Auth, handlers.DeleteHub)

	r.POST("/send/:topic", handlers.Auth, handlers.SetBroker(broker), handlers.SendMessage)           // send to all hubs
	r.POST("/hub/:slug/send/:topic", handlers.Auth, handlers.SetBroker(broker), handlers.SendMessage) // send the message only to a specific hub

	r.GET("/received/:topic/last", handlers.Auth, handlers.LastReceivedMessages)
	//r.GET("/received/:topic", handlers.Auth, handlers.ListReceivedMessages)
	//r.DELETE("/received/:topic", handlers.Auth, handlers.ClearReceivedMessages)

	log.Print("[info] Starting server on ", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
