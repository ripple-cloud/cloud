package data_test

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
)

var dbURL string

func init() {
	dbURL = os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		panic("DB_URL not set")
	}
}

func TestEncryptPassword(t *testing.T) {
	u := &data.User{}

	err := u.EncryptPassword("password")
	if err != nil {
		t.Error("EncryptPassword should not return an error for a valid password")
	}

	if u.EncryptedPassword == "" {
		t.Error("EncryptPassword should set the encrypted password in user struct")
	}
}

func TestVerifyPassword(t *testing.T) {
	u := &data.User{}

	// set an enrypted password first
	u.EncryptPassword("password")

	// returns true for correct password
	if !u.VerifyPassword("password") {
		t.Error("Expected VerifyPassword to return true")
	}

	// returns false for an incorrect password
	if u.VerifyPassword("such-password-very-secure") {
		t.Error("Expected VerifyPassword to return false")
	}
}

func TestInsert(t *testing.T) {
	// open database
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// insert new user
	u := &data.User{
		Username:          "chucknorris",
		Email:             "gmail@chucknorris.com",
		EncryptedPassword: "wood-chuck-chuck",
	}
	if err := u.Insert(db); err != nil {
		t.Error("Failed to insert user to db: %v", u)
	}

	// check if returned values are scanned back to the struct

	if u.ID == 0 {
		t.Error("ID must be set")
	}

	if u.CreatedAt == nil {
		t.Error("CreatedAt must be set")
	}

	if u.UpdatedAt == nil {
		t.Error("UpdatedAt must be set")
	}

	// TODO: Add tests for failure test paths
}

func TestGetByLogin(t *testing.T) {
	// open database
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// insert new user
	u := &data.User{
		Username:          "chucknorris1",
		Email:             "gmail1@chucknorris.com",
		EncryptedPassword: "wood-chuck-chuck",
	}
	if err := u.Insert(db); err != nil {
		t.Error("Failed to insert user to db: %v", u)
	}

	// query for the inserted user	by username
	u1 := &data.User{}
	if err := u1.GetByLogin(db, "chucknorris1"); err != nil {
		t.Error("Failed to get user with username: chucknorris1")
	}
	if u1.ID != u.ID {
		t.Error("Unexpected user record returned: %v", u1)
	}

	// query for the inserted user	by email
	u2 := &data.User{}
	if err := u2.GetByLogin(db, "gmail1@chucknorris.com"); err != nil {
		t.Error("Failed to get user with email: gmail1@chucknorris.com")
	}
	if u2.ID != u.ID {
		t.Error("Unexpected user record returned: %v", u2)
	}

	// query for a non-existing user
	u3 := &data.User{}
	err = u3.GetByLogin(db, "jackiechan")
	e, ok := err.(*data.Error)
	if !ok {
		t.Error("Returned error must be of type `data.Error`")
	}
	if e.Code != "record_not_found" {
		t.Error("Error code must be 'record_not_found' but received %s", e.Code)
	}
	if e.Desc != "user not found" {
		t.Error("Error desc must be 'user not found' but received %s", e.Desc)
	}

	// TODO: Add a test case of other errors (eg: db already closed)
}
