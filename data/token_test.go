package data_test

import (
	"testing"
	"time"

	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/testhelpers"
)

func TestTokenInsert(t *testing.T) {
	// setup database
	db := testhelpers.SetupDB(t)

	// insert new user
	u := &data.User{
		Username:          "chucknorris",
		Email:             "gmail@chucknorris.com",
		EncryptedPassword: "wood-chuck-chuck",
	}
	if err := u.Insert(db); err != nil {
		t.Error("Failed to insert user to db: %v", u)
	}

	// insert a new token for the user
	tok := &data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(),
	}
	if err := tok.Insert(db); err != nil {
		t.Error("Failed to insert token to db: %v", u)
	}

	// check if returned values are scanned back to the struct
	if tok.ID == 0 {
		t.Error("ID must be set")
	}

	if tok.CreatedAt == nil {
		t.Error("CreatedAt must be set")
	}

	// RevokedAt should be nil
	if tok.RevokedAt != nil {
		t.Error("RevokedAt must be nil")
	}

	// insert a token for non-existing user
	tok2 := &data.Token{
		UserID:    9999,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(),
	}
	if err := tok2.Insert(db); err == nil {
		t.Error("Token insert must fail if user does not exist")
	}

	db.Close()
}

func TestTokenGet(t *testing.T) {
	// setup database
	db := testhelpers.SetupDB(t)

	// insert new user
	u := &data.User{
		Username:          "chucknorris",
		Email:             "gmail@chucknorris.com",
		EncryptedPassword: "wood-chuck-chuck",
	}
	if err := u.Insert(db); err != nil {
		t.Error("Failed to insert user to db: %v", u)
	}

	// insert a new token for the user
	tok := &data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(),
	}
	if err := tok.Insert(db); err != nil {
		t.Error("Failed to insert token to db: %v", u)
	}

	// query for the inserted token
	tok1 := &data.Token{}
	if err := tok1.Get(db, tok.ID); err != nil {
		t.Error("Failed to find token for id: ", tok.ID)
	}
	if tok1.ID != tok.ID {
		t.Error("Unexpected token returned: %v", tok1)
	}

	// query for a non-existing token
	tok2 := &data.Token{}
	err := tok2.Get(db, 9999)
	if err == nil {
		t.Error("Get should return an error")
	}
	e, ok := err.(*data.Error)
	if !ok {
		t.Error("Returned error must be of type `data.Error`")
	}
	if e.Code != "record_not_found" {
		t.Error("Error code must be 'record_not_found' but received %s", e.Code)
	}
	if e.Desc != "token not found" {
		t.Error("Error desc must be 'token not found' but received %s", e.Desc)
	}

	// TODO: Add a test case of other errors (eg: db already closed)
	db.Close()
}
