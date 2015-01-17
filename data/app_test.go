package data_test

import (
	"testing"

	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/testhelpers"
)

func TestAppInsert(t *testing.T) {
	// setup database
	db := testhelpers.SetupDB(t)

	// insert new app
	a := &data.App{
		Slug:  "myapp",
		HubID: 5,
	}
	if err := a.Insert(db); err != nil {
		t.Error("Failed to insert app to db: %v", err)
	}

	// check if returned values are scanned back to the struct
	if a.ID == 0 {
		t.Error("ID must be set")
	}

	if a.CreatedAt == nil {
		t.Error("CreatedAt must be set")
	}

	if a.UpdatedAt == nil {
		t.Error("UpdatedAt must be set")
	}

	// TODO: Add tests for failure test paths
	db.Close()
}

func TestAppGet(t *testing.T) {
	// setup database
	db := testhelpers.SetupDB(t)

	a := &data.App{
		Slug:  "myapp",
		HubID: 5,
	}
	if err := a.Insert(db); err != nil {
		t.Error("Failed to insert app to db: %v", err)
	}

	// query using an existing app's slug and hub id
	a1 := &data.App{}
	if err := a1.Get(db, a.Slug, a.HubID); err != nil {
		t.Error("Failed to find an app with given slug and hub id: %v", a.Slug, a.HubID)
	}
	if a1.ID != a.ID {
		t.Error("Unexpected app record returned: %v", a1)
	}

	// query using a non-existing slug
	a2 := &data.App{}
	err := a2.Get(db, "dont-have", a.HubID)
	e, ok := err.(*data.Error)
	if !ok {
		t.Error("Returned error must be of type `data.Error`")
	}
	if e.Code != "record_not_found" {
		t.Error("Error code must be 'record_not_found' but received %s", e.Code)
	}
	if e.Desc != "app not found" {
		t.Error("Error desc must be 'app not found' but received %s", e.Desc)
	}

	// TODO: Add a test case of other errors (eg: db already closed)
	db.Close()
}

func TestAppDelete(t *testing.T) {
	// setup database
	db := testhelpers.SetupDB(t)

	a := &data.App{
		Slug:  "myapp",
		HubID: 5,
	}
	if err := a.Insert(db); err != nil {
		t.Error("Failed to insert app to db: %v", a)
	}

	// delete the created app
	if err := a.Delete(db); err != nil {
		t.Error("Failed to delete the app: %v", err)
	}

	// quering for app should return record not found
	a1 := &data.App{}
	err := a1.Get(db, a.Slug, a.HubID)
	e, ok := err.(*data.Error)
	if !ok {
		t.Error("Returned error must be of type `data.Error`")
	}
	if e.Code != "record_not_found" {
		t.Error("Error code must be 'record_not_found' but received %s", e.Code)
	}
	if e.Desc != "app not found" {
		t.Error("Error desc must be 'app not found' but received %s", e.Desc)
	}

	// TODO: Add a test case of other errors (eg: db already closed)
	db.Close()
}
