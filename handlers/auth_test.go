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

func setupServer(db *sqlx.DB, tokenSecret []byte) (*httptest.Server, error) {
	r := router.New()

	r.Default(
		handlers.SetConfig(db, []byte(tokenSecret)),
		handlers.Auth,
	)

	r.GET("/api/v0/*endpoint", func(w http.ResponseWriter, r *http.Request, c router.Context) error {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write([]byte(`{"status":"ok"}`))
		return nil
	})

	return httptest.NewServer(r), nil
}

func TestAuthToken(t *testing.T) {
	// setup DB
	db := testhelpers.SetupDB(t)
	defer db.Close()

	// setup server
	ts, err := setupServer(db, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

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

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	tCases := []testCase{
		// when access token not provided
		{"hub", http.StatusUnauthorized, `{"error":"invalid_token","error_description":"no token present in request"}`},

		// when access token is invalid
		{"hub?access_token=invalid", http.StatusUnauthorized, `{"error":"invalid_token","error_description":"token contains an invalid number of segments"}`},

		// // when access token is not properly scoped
		// // fixme currently valid scopes are ["user", "hub", "app"]
		// {"admin?access_token=" + jwt, http.statusforbidden, `{"error":"invalid_scope","error_description":"token is not valid for this scope"}`},

		// when a valid token is provided
		{"hub?access_token=" + jwt, http.StatusOK, `{"status":"ok"}`},

		// when access token is revoked
		// TODO
	}
	for _, tc := range tCases {
		res, err := http.Get(ts.URL + path.Join("/api/v0", tc.path))
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
