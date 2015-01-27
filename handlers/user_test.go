package handlers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ripple-cloud/cloud/data"
	"github.com/ripple-cloud/cloud/handlers"
	"github.com/ripple-cloud/cloud/router"
	"github.com/ripple-cloud/cloud/testhelpers"
)

func setupServerUser(db *sqlx.DB, tokenSecret []byte) (*httptest.Server, error) {
	r := router.New()

	r.Default(
		handlers.SetConfig(db, []byte(tokenSecret)),
	)

	r.GET("/signup", handlers.Signup)
	r.GET("/oauth/token", handlers.UserToken)

	return httptest.NewServer(r), nil
}

func TestSignup(t *testing.T) {
	// setup DB
	db := testhelpers.SetupDB(t)
	defer db.Close()

	// setup server
	ts, err := setupServerUser(db, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	tCases := []testCase{
		// when valid params are provided
		{"?username=foo&email=foo@example.com&password=password", http.StatusCreated, ""},

		// when username param is missing
		{"?email=foo@example.com&password=password", http.StatusBadRequest, `{"error":"username_required","error_description":"username required"}`},

		// when email param is missing
		{"?username=foo&password=password", http.StatusBadRequest, `{"error":"email_required","error_description":"email required"}`},

		// when password param is missing
		{"?username=foo&email=foo@example.com", http.StatusBadRequest, `{"error":"password_required","error_description":"password required"}`},
	}

	for i, tc := range tCases {
		res, err := http.Get(ts.URL + path.Join("/signup", tc.path))
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

		// TODO: separate successful test and parse JSON body
		// exclude testing json body for valid params because of arbitrary timestamps from postgres.
		if i != 0 {
			if body := string(b); body != tc.body {
				t.Errorf("%s - Expected response body to be %v, Got %v", tc.path, tc.body, body)
			}
		}
	}
}

func TestUserToken(t *testing.T) {
	// setup DB
	db := testhelpers.SetupDB(t)
	defer db.Close()

	// setup server
	ts, err := setupServerUser(db, []byte("secret"))
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

	tok := data.Token{
		UserID:    u.ID,
		ExpiresIn: (30 * 24 * time.Hour).Nanoseconds(), // 30 days
	}
	if err := tok.Insert(db); err != nil {
		t.Fatal(err)
	}

	//	// get the encoded JSON Web Token
	//	jwt, err := tok.EncodeJWT([]byte("secret"))
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	tCases := []testCase{
		// when valid params are provided
		// FIXME: find out why signature in jwt is different from response
		//		{"?grant_type=password&login=foo&password=password", http.StatusOK, `{"access_token":` + jwt + `","token_type":"bearer","expires_in":"720h0m0s"}`},

		// when grant_type param is invalid/missing
		{"?login=foo&password=password", http.StatusBadRequest, `{"error":"unsupported_grant_type","error_description":"supports only password grant type"}`},

		// when login param is missing
		{"?grant_type=password&password=password", http.StatusBadRequest, `{"error":"invalid_request","error_description":"login required"}`},

		// when password param is missing
		{"?grant_type=password&login=foo", http.StatusBadRequest, `{"error":"invalid_request","error_description":"password required"}`},

		// when password value is incorrect
		{"?grant_type=password&login=foo&password=abcd", http.StatusBadRequest, `{"error":"invalid_grant","error_description":"failed to authenticate user"}`},

		// when login value is incorrect
		{"?grant_type=password&login=bar&password=password", http.StatusBadRequest, `{"error":"invalid_grant","error_description":"user not found"}`},
	}

	for _, tc := range tCases {
		res, err := http.Get(ts.URL + path.Join("/oauth/token", tc.path))
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
