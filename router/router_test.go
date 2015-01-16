package router_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ripple-cloud/cloud/router"
)

func setupServer(tmpDir string) (*httptest.Server, error) {
	r := router.New()

	// directly access an endpoint
	r.GET("/endpoint", func(w http.ResponseWriter, r *http.Request, c router.Context) error {
		_, err := w.Write([]byte("OK"))
		return err
	})

	// using handlers to modify the context
	r.GET("/with_modifiers",
		func(w http.ResponseWriter, r *http.Request, c router.Context) error {
			c.Meta["test"] = "OK"
			// call next handler
			return c.Next(w, r, c)
		},
		func(w http.ResponseWriter, r *http.Request, c router.Context) error {
			o := c.Meta["test"].(string)
			_, err := w.Write([]byte(o))
			return err
		})

	// check named params and request form params
	r.GET("/with_params/:param", func(w http.ResponseWriter, r *http.Request, c router.Context) error {
		_, err := fmt.Fprintf(w, "%s %s", c.Params.ByName("param"), r.FormValue("param"))
		return err
	})

	// check for idempotency
	r.GET("/idempotent",
		func(w http.ResponseWriter, r *http.Request, c router.Context) error {
			_, err := w.Write([]byte("foo"))
			return err
		},
		func(w http.ResponseWriter, r *http.Request, c router.Context) error {
			_, err := w.Write([]byte("bar"))
			return err
		})

	// default handlers
	r.Default(func(w http.ResponseWriter, r *http.Request, c router.Context) error {
		c.Meta["default"] = "hello world"
		// call next handler
		return c.Next(w, r, c)
	})

	r.GET("/use_default_handler",
		func(w http.ResponseWriter, r *http.Request, c router.Context) error {
			o := c.Meta["default"].(string)
			_, err := w.Write([]byte(o))
			return err
		})

	// handler returning an error
	r.GET("/error",
		func(w http.ResponseWriter, r *http.Request, c router.Context) error {
			return errors.New("backend failed")
		})

	// serving static files
	err := ioutil.WriteFile(tmpDir+"/test.txt", []byte("foo bar"), 0644)
	if err != nil {
		return nil, err
	}
	r.ServeFiles("/public/*filepath", http.Dir(tmpDir))

	return httptest.NewServer(r), nil
}

func TestRouter(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "router-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	ts, err := setupServer(tmpDir)
	if err != nil {
		panic(err)
	}
	defer ts.Close()

	type testCase struct {
		path       string
		statusCode int
		body       string
	}

	// HTTP OK test cases
	tCases := []testCase{
		{"/endpoint", http.StatusOK, "OK"},
		{"/with_modifiers", http.StatusOK, "OK"},
		{"/with_params/foo?param=bar", http.StatusOK, "foo bar"},
		{"/use_default_handler", http.StatusOK, "hello world"},
		{"/error", http.StatusInternalServerError, "something went wrong\n"},
		{"/public/test.txt", http.StatusOK, "foo bar"},

		// test for idempotency (http://en.wikipedia.org/wiki/Idempotence)
		{"/idempotent", http.StatusOK, "foo"},
		{"/idempotent", http.StatusOK, "foo"},
	}

	for _, tc := range tCases {
		res, err := http.Get(ts.URL + tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != tc.statusCode {
			t.Errorf("%s - Expected status code %s, Got %s", tc.path, tc.statusCode, res.StatusCode)
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
