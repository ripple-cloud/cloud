package handlers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/handlers"
	"github.com/ripple-cloud/cloud/router"
	"github.com/ripple-cloud/cloud/testhelpers"
)

func setupServerHub(db *sqlx.DB, tokenSecret []byte) (*httptest.Server, error) {
	r := router.New()

	r.Default(
		handlers.SetConfig(db, []byte(tokenSecret)),
	)

	r.GET("/api/v0/hub", handlers.Auth, handlers.AddHub)
	r.POST("/api/v0/hub", handlers.Auth, handlers.ShowHub)
	r.DELETE("/api/v0/hub", handlers.Auth, handlers.DeleteHub)

	return httptest.NewServer(r), nil
}

func TestAddHub(t *testing.T) {
	// setup DB
	db := testhelpers.SetupDB(t)
	defer db.Close()

	// setup server
	ts, err := setupServerHub(db, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	// create a user
	u := &data.User{
		Username: "foo",
		Email:    "foo@example.com",
	}
	if err := u.EncryptPassword("password"); err != nil {
		t.Fatal(err)
	}
	if err = u.Insert(db); err != nil {
		t.Fatal(err)
	}

	// create a token for the user
	tok := data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(), // 30 days
	}
	if err := tok.Insert(db); err != nil {
		t.Fatal(err)
	}

	// get the encoded JSON Web Token
	jwt, err := tok.EncodeJWT([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	hub := data.Hub{
		Slug:   "1234",
		UserID: u.ID,
	}
	if err := hub.Insert(db); err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	jsonRespRegex := regexp.MustCompile(`^{"id":2,"slug":"abcd","user_id":1,"created_at":.+,"updated_at":.+}$`)

	// test when valid params are provided
	sc := testCase{"?slug=abcd&access_token=" + jwt, http.StatusOK, `{"id":2,"slug":"abcd","user_id":1,"created_at":.+,"updated_at":.+}`}

	res, err := http.Get(ts.URL + "/api/v0/hub" + sc.path)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != sc.statusCode {
		t.Errorf("%s - Expected status code %v, Got %v", sc.path, sc.statusCode, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if body := string(b); !jsonRespRegex.MatchString(body) {
		t.Errorf("%s - Expected response body to be %v, Got %v", sc.path, sc.body, body)
	}

	tCases := []testCase{
		// when slug param is missing
		{"?access_token=" + jwt, http.StatusBadRequest, `{"error":"invalid_request","error_description":"slug required"}`},

		// when access_token param is missing
		{"?slug=abcd", http.StatusUnauthorized, `{"error":"invalid_token","error_description":"no token present in request"}`},

		// when trying to add existing hub
		{"?slug=1234&access_token=" + jwt, http.StatusBadRequest, `{"error":"unique_violation","error_description":"hub exists"}`},
	}
	for _, tc := range tCases {
		res, err := http.Get(ts.URL + "/api/v0/hub" + tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != tc.statusCode {
			t.Errorf("%s - Expected status code %v, Got %v", tc.path, tc.statusCode, res.StatusCode)
		}
		b, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if body := string(b); body != tc.body {
			t.Errorf("%s - Expected response body to be %v, Got %v", tc.path, tc.body, body)
		}
	}
}

func TestShowHub(t *testing.T) {
	// setup DB
	db := testhelpers.SetupDB(t)
	defer db.Close()

	// setup server
	ts, err := setupServerHub(db, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	// create a user
	u := &data.User{
		Username: "foo",
		Email:    "foo@example.com",
	}
	if err := u.EncryptPassword("password"); err != nil {
		t.Fatal(err)
	}
	if err = u.Insert(db); err != nil {
		t.Fatal(err)
	}

	// create a token for the user
	tok := data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(), // 30 days
	}
	if err := tok.Insert(db); err != nil {
		t.Fatal(err)
	}

	// get the encoded JSON Web Token
	jwt, err := tok.EncodeJWT([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	hub := data.Hub{
		Slug:   "abcd",
		UserID: u.ID,
	}
	if err := hub.Insert(db); err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	tCases := []testCase{
		// when valid params are provided
		{"?access_token=" + jwt, http.StatusOK, `{"hub":["abcd"]}`},

		// when access_token param is missing
		{"?" + jwt, http.StatusUnauthorized, `{"error":"invalid_token","error_description":"no token present in request"}`},
	}
	for _, tc := range tCases {
		res, err := http.Post(ts.URL+"/api/v0/hub"+tc.path, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != tc.statusCode {
			t.Errorf("%s - Expected status code %v, Got %v", tc.path, tc.statusCode, res.StatusCode)
		}
		b, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if body := string(b); body != tc.body {
			t.Errorf("%s - Expected response body to be %v, Got %v", tc.path, tc.body, body)
		}
	}
}

func TestDeleteHub(t *testing.T) {
	// setup DB
	db := testhelpers.SetupDB(t)
	defer db.Close()

	// setup server
	ts, err := setupServerHub(db, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	// create a user
	u := &data.User{
		Username: "foo",
		Email:    "foo@example.com",
	}
	if err := u.EncryptPassword("password"); err != nil {
		t.Fatal(err)
	}
	if err = u.Insert(db); err != nil {
		t.Fatal(err)
	}

	// create a token for the user
	tok := data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(), // 30 days
	}
	if err := tok.Insert(db); err != nil {
		t.Fatal(err)
	}

	// get the encoded JSON Web Token
	jwt, err := tok.EncodeJWT([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	hub := data.Hub{
		Slug:   "abcd",
		UserID: u.ID,
	}
	if err := hub.Insert(db); err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	jsonRespRegex := regexp.MustCompile(`^{"id":1,"slug":"abcd","user_id":1,"created_at":.+,"updated_at":.+}$`)

	// test when valid params are provided
	sc := testCase{"?slug=abcd&access_token=" + jwt, http.StatusOK, `{"id":1,"slug":"abcd","user_id":1,"created_at":.+,"updated_at":.+}`}
	result, err := http.NewRequest("DELETE", ts.URL+"/api/v0/hub"+sc.path, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, _ := http.DefaultClient.Do(result)

	if res.StatusCode != sc.statusCode {
		t.Errorf("%s - Expected status code %v, Got %v", sc.path, sc.statusCode, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if body := string(b); !jsonRespRegex.MatchString(body) {
		t.Errorf("%s - Expected response body to be %v, Got %v", sc.path, sc.body, body)
	}

	tCases := []testCase{
		// when slug param is missing
		{"?access_token=" + jwt, http.StatusBadRequest, `{"error":"invalid_request","error_description":"slug required"}`},

		// when hub does not exist
		{"?slug=1234&access_token=" + jwt, http.StatusBadRequest, `{"error":"record_not_found","error_description":"hub not found"}`},
	}
	for _, tc := range tCases {
		result, err := http.NewRequest("DELETE", ts.URL+"/api/v0/hub"+tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, _ := http.DefaultClient.Do(result)

		if res.StatusCode != tc.statusCode {
			t.Errorf("%s - Expected status code %v, Got %v", tc.path, tc.statusCode, res.StatusCode)
		}
		b, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if body := string(b); body != tc.body {
			t.Errorf("%s - Expected response body to be %v, Got %v", tc.path, tc.body, body)
		}
	}
}
