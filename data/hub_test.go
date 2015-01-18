package data_test

import (
	"testing"

	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/testhelpers"
)

func TestHubInsert(t *testing.T) {
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

	h := &data.Hub{
		Slug:   "earthworm",
		UserID: 1,
	}
	if err := h.Insert(db); err != nil {
		t.Error("Failed to insert h to db: %v", h)
	}

	// check if returned values are scanned back to the struct
	if h.ID == 0 {
		t.Error("ID must be set")
	}

	if h.CreatedAt == nil {
		t.Error("CreatedAt must be set")
	}

	if h.UpdatedAt == nil {
		t.Error("UpdatedAt must be set")
	}

	// TODO: Add tests for failure test paths
	db.Close()
}

func TestHubGetByUserId(t *testing.T) {
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

	h := &data.Hub{
		Slug:   "earthworm",
		UserID: 1,
	}
	if err := h.Insert(db); err != nil {
		t.Error("Failed to insert h to db: %v", h)
	}

	// query for the inserted hub by userid
	var h1 data.Hubs
	if err := h1.GetByUserId(db, 1); err != nil {
		t.Error("Failed to get hubs with userid: 1")
	}
	if h1[0] != h.Slug {
		t.Error("Unexpected user record returned: %v", h1)
	}

	// query for a non-existing hub
	var h2 data.Hubs
	err := h2.GetByUserId(db, 9999)
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
	if e.Desc != "hub not found" {
		t.Error("Error desc must be 'user not found' but received %s", e.Desc)
	}

	// TODO: Add a test case of other errors (eg: db already closed)
	db.Close()
}

func TestHubGetByHub(t *testing.T) {
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

	h := &data.Hub{
		Slug:   "earthworm",
		UserID: 1,
	}
	if err := h.Insert(db); err != nil {
		t.Error("Failed to insert h to db: %v", h)
	}

	// query for the inserted hub by hub
	h1 := &data.Hub{}
	if err := h1.GetByHub(db, "earthworm"); err != nil {
		t.Error("Failed to get hubs with hub: earthworm")
	}
	if h1.UserID != h.UserID {
		t.Error("Unexpected user record returned: %v", h1)
	}

	// query for a non-existing hub
	h2 := &data.Hub{}
	err := h2.GetByHub(db, "snail")
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
	if e.Desc != "hub not found" {
		t.Error("Error desc must be 'user not found' but received %s", e.Desc)
	}

	//TODO: Add a test case of other errors (eg: db already closed)
	db.Close()
}

func TestHubDelete(t *testing.T) {
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

	h := &data.Hub{
		Slug:   "earthworm",
		UserID: 1,
	}
	if err := h.Insert(db); err != nil {
		t.Error("Failed to insert h to db: %v", h)
	}

	// query for the inserted hub by hub
	h1 := &data.Hub{
		Slug: "earthworm",
	}
	if err := h1.Delete(db); err != nil {
		t.Error("Failed to get hubs with hub: earthworm")
	}
	if h1.Slug != h.Slug {
		t.Error("Unexpected user record returned: %v", h1)
	}

	// check if returned values from Delete are scanned back to the struct
	if h1.ID == 0 {
		t.Error("ID must be set")
	}

	if h1.CreatedAt == nil {
		t.Error("CreatedAt must be set")
	}

	if h1.UpdatedAt == nil {
		t.Error("UpdatedAt must be set")
	}

	// query for a non-existing hub
	h2 := &data.Hub{
		Slug: "snail",
	}
	err := h2.Delete(db)
	if err == nil {
		t.Error("Delete should return an error")
	}

	//TODO: Add a test case of other errors (eg: db already closed)
	db.Close()
}
