package data_test

import (
	"testing"

	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/testhelpers"
)

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

	// check if adding existing hub violates unique constraint
	u1 := &data.User{
		Username:          "chucknorris",
		Email:             "gmail@chucknorris.com",
		EncryptedPassword: "wood-chuck-chuck",
	}
	err := u1.Insert(db)

	if err == nil {
		t.Error("Insert should return an error")
	}
	e, ok := err.(*data.Error)
	if !ok {
		t.Error("Returned error must be of type `data.Error`")
	}
	if e.Code != "unique_violation" {
		t.Error("Error code must be 'unique violation' but received %s", e.Code)
	}
	if e.Desc != "user exists" {
		t.Error("Error desc must be 'user exists' but received %s", e.Desc)
	}

	db.Close()
}

func TestGetByLogin(t *testing.T) {
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

	// query for the inserted user by username
	u1 := &data.User{}
	if err := u1.GetByLogin(db, "chucknorris"); err != nil {
		t.Error("Failed to get user with username: chucknorris")
	}
	if u1.ID != u.ID {
		t.Error("Unexpected user record returned: %v", u1)
	}

	// query for the inserted user by email
	u2 := &data.User{}
	if err := u2.GetByLogin(db, "gmail@chucknorris.com"); err != nil {
		t.Error("Failed to get user with email: gmail@chucknorris.com")
	}
	if u2.ID != u.ID {
		t.Error("Unexpected user record returned: %v", u2)
	}

	// query for a non-existing user
	u3 := &data.User{}
	err := u3.GetByLogin(db, "jackiechan")
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
	db.Close()
}

func TestGet(t *testing.T) {
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

	// query using an existing user ID
	u1 := &data.User{}
	if err := u1.Get(db, u.ID); err != nil {
		t.Error("Failed to find a user with given ID: %v", u.ID)
	}
	if u1.ID != u.ID {
		t.Error("Unexpected user record returned: %v", u1)
	}

	// query using a non-existing user ID
	u2 := &data.User{}
	err := u2.Get(db, 9999)
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
	db.Close()
}
